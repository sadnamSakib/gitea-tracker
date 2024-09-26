package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/config"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/db"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ClearUsers() error {
	collection := db.MongoDatabase.Collection(userCollection)
	err := collection.Drop(context.Background())
	if err != nil {
		return err
	}
	return nil
}
func FollowUser(userName string) error {
	collection := db.MongoDatabase.Collection(userCollection)
	filter := bson.M{"username": userName}
	update := bson.M{
		"$set": bson.M{
			"following": true,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func UnfollowUser(userName string) error {
	collection := db.MongoDatabase.Collection(userCollection)
	filter := bson.M{"username": userName}
	update := bson.M{
		"$set": bson.M{
			"following": false,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return err
	}
	return nil

}

func GetAllUsers(page, limit string) ([]model.User, error) {
	users := make([]model.User, 0)
	collection := db.MongoDatabase.Collection(userCollection)
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
	cursor, err := collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		return nil, err
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

func GetUser(userName string) (model.User, error) {
	user := model.User{}
	collection := db.MongoDatabase.Collection(userCollection)
	filter := bson.M{"username": userName}
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func SearchUsers(query, page, limit string) ([]model.User, error) {
	users := make([]model.User, 0)
	collection := db.MongoDatabase.Collection(userCollection)
	filter := bson.M{
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
		return nil, err
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

func FetchUsersFromGitea(page int) ([]model.User, error) {
	users := make([]model.User, 0)

	url := fmt.Sprintf("%s/admin/users?page=%d&access_token=%s", config.AppConfig.GITEA.Base_URL, page, config.AppConfig.GITEA.API_KEY)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return users, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return users, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return users, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return users, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if len(users) == 0 {
		return users, nil
	}

	next_users, err := FetchUsersFromGitea(page + 1)
	if err != nil {
		return users, err
	}
	users = append(users, next_users...)
	return users, nil
}
func SyncUsersWithDB(users []model.User) error {
	collection := db.MongoDatabase.Collection(userCollection)
	existingUsers, err := GetAllUsers("", "")
	if err != nil {
		return err
	}
	existingUserMap := make(map[string]model.User)
	for _, user := range users {
		existingUserMap[user.Username] = user
	}
	documentsToBeAdded := make([]interface{}, 0, len(users))

	for _, user := range users {

		filter := bson.M{"username": user.Username}
		var existingUser model.User
		err := collection.FindOne(context.Background(), filter).Decode(&existingUser)

		if err == mongo.ErrNoDocuments {
			documentsToBeAdded = append(documentsToBeAdded, user)
			continue
		}
	}
	for _, user := range existingUsers {
		if _, ok := existingUserMap[user.Username]; !ok {
			_, err := collection.DeleteOne(context.Background(), bson.M{"username": user.Username})
			if err != nil {
				fmt.Println("Error deleting user ", user.Username)
				continue
			}

		}
	}

	if len(documentsToBeAdded) == 0 {
		return nil
	}
	_, err = collection.InsertMany(context.Background(), documentsToBeAdded)
	if err != nil {
		return err
	}

	return nil
}
