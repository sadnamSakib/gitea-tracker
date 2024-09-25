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
	documents := make([]interface{}, len(repos))
	for i, repo := range repos {
		documents[i] = repo
	}

	_, err := collection.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	return nil
}
