package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/config"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/db"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
)

const orgCollection = "orgs"
const repoCollection = "repos"

func GetAllOrganizationsFromGitea(page int, orgs *[]*model.Org) error {
	// Define the Gitea API URL with access token
	currentPageOrgs := []model.Org{}
	limit := 50
	url := fmt.Sprintf("https://gitea.vivasoftltd.com/api/v1/orgs?page=%d&limit=%d&access_token=%s", page, limit, config.AppConfig.GITEA.API_KEY)

	// Create a new HTTP client and request
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

	if err := json.NewDecoder(resp.Body).Decode(&currentPageOrgs); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if len(currentPageOrgs) == 0 {
		return nil
	}
	for _, org := range currentPageOrgs {
		*orgs = append(*orgs, &org)
	}

	fmt.Println(len(*orgs))

	return nil
}
func GetAllOrganizationFromDB(orgs *[]*model.Org) error {
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
func SyncOrganizations(orgs []*model.Org) error {

	collection := db.MongoDatabase.Collection(orgCollection)
	err := collection.Drop(context.Background())
	if err != nil {
		return err
	}

	documents := make([]interface{}, len(orgs))
	for i, org := range orgs {
		documents[i] = org
	}

	_, err = collection.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	return nil
}

func GetAllRepoOfOrganization(page int, orgName string, repos *[]*model.Repo) error {
	// Define the Gitea API URL with access token
	currentPageRepos := []model.Repo{}

	url := fmt.Sprintf("https://gitea.vivasoftltd.com/api/v1/orgs/%s/repos?page=%d&access_token=%s", orgName, page, config.AppConfig.GITEA.API_KEY)

	// Create a new HTTP client and request
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

	if err := json.NewDecoder(resp.Body).Decode(&currentPageRepos); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if len(currentPageRepos) == 0 {
		return nil
	}
	for _, org := range currentPageRepos {
		*repos = append(*repos, &org)
	}

	GetAllRepoOfOrganization(page+1, orgName, repos)

	return nil
}

func SyncRepos(repos []*model.Repo) error {
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

func GetAllRepoFromDB(orgName string, repos *[]*model.Repo) error {
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
