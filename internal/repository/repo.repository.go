package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/config"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/db"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ClearRepos() error {
	collection := db.MongoDatabase.Collection(repoCollection)
	err := collection.Drop(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func SearchUsersOfRepo(org, repo, query, page, limit string) ([]model.User, error) {
	users := make([]model.User, 0)
	collection := db.MongoDatabase.Collection(userCollection)
	filter := bson.M{
		"repos": bson.M{
			"$elemMatch": bson.M{
				"name":           repo,
				"owner.username": org,
			},
		},
		"username": bson.M{
			"$regex":   query,
			"$options": "i",
		},
	}
	findOptions := options.Find()
	if page != "" {
		pageNum, err := strconv.Atoi(page)
		if err != nil {
			return nil, fmt.Errorf("invalid page number: %w", err)
		}
		limitNum := 10
		if limit != "" {
			limitNum, err = strconv.Atoi(limit)
			if err != nil {
				return nil, fmt.Errorf("invalid limit number: %w", err)
			}
		}

		findOptions.SetLimit(int64(limitNum))
		findOptions.SetSkip(int64((pageNum - 1) * limitNum))
	}
	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
func FetchRepoOfOrgFromGitea(page int, orgName string) ([]model.Repo, error) {

	repos := make([]model.Repo, 0)

	url := fmt.Sprintf("%s/orgs/%s/repos?page=%d&access_token=%s", config.AppConfig.GITEA.Base_URL, orgName, page, config.AppConfig.GITEA.API_KEY)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return repos, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return repos, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return repos, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return repos, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if len(repos) == 0 {
		return []model.Repo{}, nil
	}
	next_repos, err := FetchRepoOfOrgFromGitea(page+1, orgName)
	if err != nil {
		return repos, err
	}
	repos = append(repos, next_repos...)
	return repos, nil
}

func SyncReposWithDB(repos []model.Repo) error {
	collection := db.MongoDatabase.Collection(repoCollection)

	var reposToInsert []interface{}

	for _, repo := range repos {

		filter := bson.M{
			"name":           repo.Name,
			"owner.username": repo.Owner.Username,
		}

		err := collection.FindOne(context.Background(), filter).Err()
		if err == mongo.ErrNoDocuments {

			reposToInsert = append(reposToInsert, repo)
		} else if err != nil {

			return fmt.Errorf("failed to check existing repository: %v", err)
		}
	}

	if len(reposToInsert) > 0 {
		_, err := collection.InsertMany(context.Background(), reposToInsert)
		if err != nil {
			return fmt.Errorf("failed to insert new repositories: %v", err)
		}
	}

	return nil
}

func GetRepoActivityByDateRange(repoName string, start_date_str string, end_date_str string) ([]model.Activity, error) {
	collection := db.MongoDatabase.Collection(activitesCollection)
	layout := "2006-01-02"

	filter := bson.M{
		"repo.name": repoName,
	}

	if start_date_str != "" {
		start_date, err := time.Parse(layout, start_date_str)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format: %v", err)
		}
		filter["date"] = bson.M{"$gte": start_date}
	}

	if end_date_str != "" {
		end_date, err := time.Parse(layout, end_date_str)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %v", err)
		}

		end_date = end_date.Add(24*time.Hour - time.Nanosecond)
		if _, exists := filter["date"]; exists {
			filter["date"].(bson.M)["$lte"] = end_date
		} else {
			filter["date"] = bson.M{"$lte": end_date}
		}
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
func SyncRepoActivitiesWithDB(repoName string, activities []model.Activity) error {

	if len(activities) == 0 {

		return nil
	}

	collection := db.MongoDatabase.Collection(repoCollection)
	filter := bson.M{"name": repoName}

	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	lastMonday := now.AddDate(0, 0, -weekday+1).Format("2006-01-02")
	lastMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	lastYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	weeklyCommits, monthlyCommits, yearlyCommits, allTimeCommits := AggregateCommits(lastMonday, lastMonth, lastYear, activities)

	update := bson.M{
		"$set": bson.M{
			"aggregated_commits.last_week":  weeklyCommits,
			"aggregated_commits.last_month": monthlyCommits,
			"aggregated_commits.last_year":  yearlyCommits,
			"aggregated_commits.all_time":   allTimeCommits,
		},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
func SyncNewRepoActivitiesWithDB(repoName string, activities []model.Activity) error {

	if len(activities) == 0 {

		return nil
	}

	collection := db.MongoDatabase.Collection(repoCollection)
	filter := bson.M{"name": repoName}

	var repo model.Repo
	err := collection.FindOne(context.Background(), filter).Decode(&repo)
	if err != nil {
		return err
	}

	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	lastMonday := now.AddDate(0, 0, -weekday+1).Format("2006-01-02")
	lastMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	lastYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	weeklyCommits, monthlyCommits, yearlyCommits, allTimeCommits := AggregateCommits(lastMonday, lastMonth, lastYear, activities)

	update := bson.M{
		"$inc": bson.M{
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

func GetNewRepoActivity(repo string) ([]model.Activity, error) {
	collection := db.MongoDatabase.Collection(systemCollection)
	system := collection.FindOne(context.Background(), bson.M{"name": "system"})
	if system.Err() != nil {
		return nil, system.Err()
	}
	var systemDoc model.SystemSummary
	err := system.Decode(&systemDoc)
	if err != nil {
		return nil, err
	}
	lastSync := systemDoc.Last_synced
	collection = db.MongoDatabase.Collection(activitesCollection)
	filter := bson.M{"name": repo, "date": bson.M{"$gt": lastSync}}
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
