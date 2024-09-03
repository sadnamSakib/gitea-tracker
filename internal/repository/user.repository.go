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

const userCollection = "users"
const activitesCollection = "activities"

func GetAllUserFromGitea(page int, orgName string, users *[]*model.User) error {
	currentPageUsers := []model.User{}
	limit := 50
	url := fmt.Sprintf("https://gitea.vivasoftltd.com/api/v1/orgs/%s/members?page=%d&limit=%d&access_token=%s", orgName, page, limit, config.AppConfig.GITEA.API_KEY)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&currentPageUsers); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	currentTime := time.Now()

	if len(currentPageUsers) == 0 {
		return nil
	}
	for _, user := range currentPageUsers {
		user.Last_updated = currentTime
		*users = append(*users, &user)
	}
	GetAllUserFromGitea(page+1, orgName, users)

	return nil
}

func GetAllUserActivityFromGitea(page int, userName string, activities *[]*model.Activity) error {
	currentPageActivities := []model.Activity{}
	limit := 50
	url := fmt.Sprintf("https://gitea.vivasoftltd.com/api/v1/users/%s/activities/feeds?only-performed-by=true&page=%d&limit=%d&access_token=%s", userName, page, limit, config.AppConfig.GITEA.API_KEY)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&currentPageActivities); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if len(currentPageActivities) == 0 {
		return nil
	}
	for _, activity := range currentPageActivities {
		if activity.OpType == "commit_repo" {
			*activities = append(*activities, &activity)
		}

	}
	GetAllUserActivityFromGitea(page+1, userName, activities)

	return nil
}

func SyncUsers(users []*model.User) error {
	collection := db.MongoDatabase.Collection(userCollection)
	err := collection.Drop(context.Background())
	if err != nil {
		return err
	}

	uniqueUsers := make(map[int]*model.User)
	for _, user := range users {
		if _, exists := uniqueUsers[user.Id]; !exists {
			uniqueUsers[user.Id] = user
		}
	}
	for _, user := range uniqueUsers {
		fmt.Println(user.Username)
	}
	fmt.Println(len(uniqueUsers))

	documents := make([]interface{}, 0, len(uniqueUsers))

	for _, user := range uniqueUsers {
		documents = append(documents, user)
	}
	fmt.Println(len(documents))

	_, err = collection.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	return nil
}

func SyncActivities(activities []*model.Activity) error {
	collection := db.MongoDatabase.Collection(activitesCollection)

	documents := make([]interface{}, len(activities))
	for i, user := range activities {
		documents[i] = user
	}

	if len(documents) == 0 {
		return nil
	}

	_, err := collection.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	username := activities[0].PerformedBy.Username
	collection = db.MongoDatabase.Collection(userCollection)
	filter := bson.M{"username": username}

	update := bson.M{
		"$set": bson.M{
			"last_updated": time.Now(),
			"TotalCommits": len(activities),
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func GetAllUsersFromDB(users *[]*model.User) error {
	collection := db.MongoDatabase.Collection(userCollection)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			return err
		}
		*users = append(*users, &user)
	}
	return nil

}

func GetUser(userName string, user *model.User) error {
	collection := db.MongoDatabase.Collection(userCollection)
	filter := bson.M{"username": userName}
	err := collection.FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		return err
	}
	return nil
}

func GetUserActivityByDateRange(userName string, start_date_str string, end_date_str string, repo string) ([]*model.Activity, error) {
	collection := db.MongoDatabase.Collection(activitesCollection)
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

	// If end_date_str is not empty, parse it and add to the filter
	if end_date_str != "" {
		end_date, err := time.Parse(layout, end_date_str)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %v", err)
		}
		// To ensure the end date includes the entire day
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
	var activities []*model.Activity
	for cursor.Next(context.Background()) {
		var activity model.Activity
		if err := cursor.Decode(&activity); err != nil {
			return nil, err
		}
		activities = append(activities, &activity)
	}
	return activities, nil
}

func GetDailyUserActivity(page int, userName string, date string, activities *[]*model.Activity) error {
	currentPageActivities := []model.Activity{}
	limit := 50
	url := fmt.Sprintf("https://gitea.vivasoftltd.com/api/v1/users/%s/activities/feeds?only-performed-by=true&page=%d&limit=%d&date=%s&access_token=%s", userName, page, limit, date, config.AppConfig.GITEA.API_KEY)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&currentPageActivities); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if len(currentPageActivities) == 0 {
		return nil
	}
	for _, activity := range currentPageActivities {
		if activity.OpType == "commit_repo" {
			*activities = append(*activities, &activity)
		}

	}
	GetAllUserActivityFromGitea(page+1, userName, activities)

	return nil
}

func SyncDailyActivities(activities []*model.Activity) error {
	collection := db.MongoDatabase.Collection(activitesCollection)

	documents := make([]interface{}, len(activities))
	for i, user := range activities {
		documents[i] = user
	}

	if len(documents) == 0 {
		return nil
	}

	_, err := collection.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	username := activities[0].PerformedBy.Username
	collection = db.MongoDatabase.Collection(userCollection)
	filter := bson.M{"username": username}

	update := bson.M{
		"$set": bson.M{
			"last_updated": time.Now(),
			"TotalCommits": bson.M{"$sum": len(activities)},
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
