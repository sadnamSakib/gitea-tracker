package repository

import (
	"context"
	"fmt"
	"time"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/db"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
)

const userCollection = "users"
const activitesCollection = "activities"

func GetAllUsers(users *[]*model.User) error {
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
