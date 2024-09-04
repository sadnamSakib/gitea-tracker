package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/db"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userCollection = "users"
const activitesCollection = "activities"

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
	err := collection.FindOne(context.Background(), filter).Decode(user)
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

func GetUserActivityByDateRange(userName string, start_date_str string, end_date_str string, repo string) ([]model.Activity, error) {
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
