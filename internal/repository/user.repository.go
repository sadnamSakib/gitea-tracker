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
		user.LastUpdated = currentTime
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

	documents := make([]interface{}, len(users))
	for i, user := range users {
		documents[i] = user
	}

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

	_, err := collection.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	username := activities[0].PerformedBy.Username
	collection = db.MongoDatabase.Collection(userCollection)
	filter := bson.M{"username": username}

	// Define the update to set both last_updated and TotalCommits
	update := bson.M{
		"$set": bson.M{
			"last_updated": time.Now(),
			"TotalCommits": len(activities),
		},
	}

	// Perform the update
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
