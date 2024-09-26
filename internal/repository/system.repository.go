package repository

import (
	"context"
	"time"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/db"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SyncSystemSummary(orgs, repos, users int, last_synced time.Time, elapsedTime int64) error {
	collection := db.MongoDatabase.Collection(systemCollection)
	filter := bson.M{"username": "system"}
	update := bson.M{
		"$set": bson.M{
			"total_orgs":            orgs,
			"total_repos":           repos,
			"total_users":           users,
			"last_synced":           last_synced,
			"is_synced":             true,
			"estimated_update_time": elapsedTime,
		},
	}

	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func GetSystemSummary() (model.SystemSummary, error) {
	collection := db.MongoDatabase.Collection(systemCollection)
	filter := bson.M{"username": "system"}
	var summary model.SystemSummary
	err := collection.FindOne(context.Background(), filter).Decode(&summary)
	if err != nil {
		return model.SystemSummary{}, err
	}
	return summary, nil
}

func UpdateSyncStatus(is_synced bool) error {
	collection := db.MongoDatabase.Collection(systemCollection)
	filter := bson.M{"username": "system"}
	update := bson.M{
		"$set": bson.M{
			"is_synced": is_synced,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}
