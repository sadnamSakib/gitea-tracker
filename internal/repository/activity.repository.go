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

func GetUserActivityByDateRange(userName string, start_date_str string, end_date_str string, repo string) ([]model.Activity, error) {
	collection := db.MongoDatabase.Collection(activitesCollection)
	fmt.Println(start_date_str)
	fmt.Println(end_date_str)
	layout := "2006-01-02"
	filter := bson.M{
		"performedby.username": userName,
	}
	if start_date_str != "" {
		start_date, err := time.Parse(layout, start_date_str)
		if err != nil {
			return nil, err
		}

		filter["date"] = bson.M{"$gte": start_date}
	}

	if end_date_str != "" {
		end_date, err := time.Parse(layout, end_date_str)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %v", err)
		}

		end_date = end_date.Add(24*time.Hour - time.Nanosecond)
		if filter["date"] != nil {
			filter["date"].(bson.M)["$lte"] = end_date
		} else {
			filter["date"] = bson.M{"$lte": end_date}
		}

	}
	if repo != "" {
		filter["repo.name"] = repo
	}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var activities []model.Activity
	for cursor.Next(context.Background()) {
		var activity model.Activity
		if err := cursor.Decode(&activity); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, nil
}
func FetchUserActivityFromGitea(page int, userName string) ([]model.Activity, error) {
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

	if len(activities) == 0 {
		return activities, nil
	}
	commitActivities := make([]model.Activity, 0)
	for _, activity := range activities {
		if activity.OpType == "commit_repo" {
			commitActivities = append(commitActivities, activity)
		}
	}

	next_activities, err := FetchUserActivityFromGitea(page+1, userName)
	if err != nil {
		return activities, err
	}
	commitActivities = append(commitActivities, next_activities...)
	return commitActivities, nil
}

func SyncActivitiesWithDB(username string, activities []model.Activity) error {
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
	repos := make(map[string]model.Repo)
	for i, user := range activities {
		repos[user.Repo.Name] = user.Repo
		documents[i] = user
	}
	_, err := collection.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	collection = db.MongoDatabase.Collection(userCollection)
	filter := bson.M{"username": username}
	repoList := make([]model.Repo, 0, len(repos))
	for _, repo := range repos {
		repoList = append(repoList, repo)
	}

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
			"last_updated":                  lastUpdateTime,
			"total_commits":                 len(activities),
			"repos":                         repoList,
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
func ClearActivities() error {
	collection := db.MongoDatabase.Collection(activitesCollection)
	err := collection.Drop(context.Background())
	if err != nil {
		return err
	}
	return nil
}
