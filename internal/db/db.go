package db

import (
	"context"
	"log"
	"time"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

func Connect() {
	dbConfig := config.AppConfig.Database.MongoDB

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConfig.URI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %s", err)
	}

	log.Println("Connected to MongoDB!")

	MongoClient = client
	MongoDatabase = client.Database(dbConfig.Database)
}

func Disconnect() {

	if err := MongoClient.Disconnect(context.Background()); err != nil {
		log.Fatalf("Failed to disconnect MongoDB client: %s", err)
	}
	log.Println("Disconnected from MongoDB.")
}
