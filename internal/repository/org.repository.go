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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ClearOrgs() error {
	collection := db.MongoDatabase.Collection(orgCollection)
	err := collection.Drop(context.Background())
	if err != nil {
		return err
	}
	return nil
}
func SyncOrgsWithDB(orgs []model.Org) error {

	collection := db.MongoDatabase.Collection(orgCollection)

	documents := make([]interface{}, len(orgs))
	for i, org := range orgs {
		documents[i] = org
	}

	_, err := collection.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	return nil
}

func GetAllOrgs() ([]*model.Org, error) {
	orgs := make([]*model.Org, 0)
	collection := db.MongoDatabase.Collection(orgCollection)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {

		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var org model.Org
		if err = cursor.Decode(&org); err != nil {

			return nil, err
		}
		orgs = append(orgs, &org)
	}
	return orgs, nil

}
func FetchOrgsFromGitea(page int) ([]model.Org, error) {
	orgs := make([]model.Org, 0)
	url := fmt.Sprintf("%s/orgs?page=%d&access_token=%s", config.AppConfig.GITEA.Base_URL, page, config.AppConfig.GITEA.API_KEY)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return orgs, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return orgs, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return orgs, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&orgs); err != nil {
		return orgs, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if len(orgs) == 0 {
		return []model.Org{}, nil
	}
	next_orgs, err := FetchOrgsFromGitea(page + 1)
	if err != nil {
		return orgs, err
	}
	orgs = append(orgs, next_orgs...)
	return orgs, nil
}

func GetAllReposFromOrg(orgName, page, limit string) ([]model.Repo, error) {
	repos := make([]model.Repo, 0)
	collection := db.MongoDatabase.Collection(repoCollection)
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
	cursor, err := collection.Find(context.Background(), bson.M{"owner.username": orgName}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var repo model.Repo
		if err = cursor.Decode(&repo); err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

func GetRepo(orgName, repoName string) (model.Repo, error) {
	repo := model.Repo{}
	collection := db.MongoDatabase.Collection(repoCollection)
	err := collection.FindOne(context.Background(), bson.M{"owner.username": orgName, "name": repoName}).Decode(&repo)
	if err != nil {

		return model.Repo{}, err
	}
	return repo, nil
}

func GetAllUsersFromRepo(org, repo, page, limit string) ([]model.User, error) {

	collection := db.MongoDatabase.Collection("users")
	users := make([]model.User, 0)
	filter := bson.M{
		"repos": bson.M{
			"$elemMatch": bson.M{
				"name":           repo,
				"owner.username": org,
			},
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
		if err = cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}
		users = append(users, user)
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor encountered error: %w", err)
	}

	return users, nil
}

func SearchRepos(org, query, page, limit string) ([]model.Repo, error) {
	repos := make([]model.Repo, 0)
	collection := db.MongoDatabase.Collection(repoCollection)
	filter := bson.M{
		"owner.username": org,
		"name": bson.M{
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
		var repo model.Repo
		if err = cursor.Decode(&repo); err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}
	return repos, nil
}
