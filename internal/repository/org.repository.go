package repository

import (
	"context"
	"fmt"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/db"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
)

const orgCollection = "orgs"
const repoCollection = "repos"

func GetAllOrgs(orgs *[]*model.Org) error {
	collection := db.MongoDatabase.Collection(orgCollection)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {

		return err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var org model.Org
		if err = cursor.Decode(&org); err != nil {

			return err
		}
		*orgs = append(*orgs, &org)
	}
	return nil

}

func GetAllReposFromOrg(orgName string, repos *[]*model.Repo) error {
	collection := db.MongoDatabase.Collection(repoCollection)
	cursor, err := collection.Find(context.Background(), bson.M{"owner.username": orgName})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var repo model.Repo
		if err = cursor.Decode(&repo); err != nil {
			return err
		}
		*repos = append(*repos, &repo)
	}
	return nil
}

func GetAllUsersFromRepo(org string, repo string, users *[]*model.User) error {
	collection := db.MongoDatabase.Collection("users")
	filter := bson.M{
		"repos": bson.M{
			"$elemMatch": bson.M{
				"name":           repo,
				"owner.username": org,
			},
		},
	}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(context.Background())

	// Iterate over the cursor to decode each user and append to the users slice
	for cursor.Next(context.Background()) {
		var user model.User
		if err = cursor.Decode(&user); err != nil {
			return fmt.Errorf("failed to decode user: %w", err)
		}
		*users = append(*users, &user)
	}

	// Check for any errors that occurred during iteration
	if err = cursor.Err(); err != nil {
		return fmt.Errorf("cursor encountered error: %w", err)
	}

	return nil
}
