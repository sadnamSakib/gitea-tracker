package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/config"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/db"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
)

func FetchNewUserActivityFromGitea(page int, userName string, lastUpdateTime time.Time) ([]model.Activity, error) {

	activities := make([]model.Activity, 0)
	url := fmt.Sprintf("%s/users/%s/activities/feeds?only-performed-by=true&page=%d&access_token=%s", config.AppConfig.GITEA.Base_URL, userName, page, config.AppConfig.GITEA.API_KEY)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return activities, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return activities, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return activities, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&activities); err != nil {
		return activities, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	commitActivities := make([]model.Activity, 0)
	for _, activity := range activities {
		if activity.OpType == "commit_repo" && activity.Date.After(lastUpdateTime) {
			fmt.Println(userName, activity.Date, lastUpdateTime)
			commitActivities = append(commitActivities, activity)
		}
		if activity.Date.Before(lastUpdateTime) {
			return commitActivities, nil
		}
	}

	if len(commitActivities) == 0 {
		return commitActivities, nil
	}
	next_activities, err := FetchNewUserActivityFromGitea(page+1, userName, lastUpdateTime)
	if err != nil {
		return activities, err
	}
	commitActivities = append(commitActivities, next_activities...)
	return commitActivities, nil
}

func SyncNewActivitiesWithDB(username string, activities []model.Activity) error {
	if len(activities) == 0 {
		collection := db.MongoDatabase.Collection(userCollection)
		filter := bson.M{"username": username}

		update := bson.M{
			"$set": bson.M{
				"last_updated": time.Now().UTC(),
			},
		}

		_, err := collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return err
		}
		return nil
	}

	collection := db.MongoDatabase.Collection(activitesCollection)
	documents := make([]interface{}, len(activities))
	reposSet := make(map[string]model.Repo)
	for i, activity := range activities {
		documents[i] = activity
		reposSet[activity.Repo.Name] = activity.Repo
	}
	_, err := collection.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	collection = db.MongoDatabase.Collection(userCollection)
	filter := bson.M{"username": username}

	var user model.User
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return err
	}

	for _, repo := range user.Repos {
		reposSet[repo.Name] = repo
	}

	repoList := make([]model.Repo, 0, len(reposSet))
	for _, repo := range reposSet {
		repoList = append(repoList, repo)
	}
	fmt.Println("Total repos : ", len(reposSet), " for user ", username)
	lastUpdateTime := time.Now().UTC()
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	lastMonday := now.AddDate(0, 0, -weekday+1).Format("2006-01-02")
	lastMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	lastYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	weeklyCommits, monthlyCommits, yearlyCommits, allTimeCommits := AggregateUserCommits(lastMonday, lastMonth, lastYear, activities)

	update := bson.M{
		"$set": bson.M{
			"last_updated": lastUpdateTime,
			"repos":        repoList,
		},
		"$inc": bson.M{
			"total_commits":                 len(activities),
			"aggregated_commits.last_week":  weeklyCommits,
			"aggregated_commits.last_month": monthlyCommits,
			"aggregated_commits.last_year":  yearlyCommits,
			"aggregated_commits.all_time":   allTimeCommits,
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
