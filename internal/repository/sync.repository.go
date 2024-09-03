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

func FetchOrgsFromGitea(page int, orgs *[]*model.Org) error {
	currentPageOrgs := []model.Org{}
	limit := 50
	url := fmt.Sprintf("https://gitea.vivasoftltd.com/api/v1/orgs?page=%d&limit=%d&access_token=%s", page, limit, config.AppConfig.GITEA.API_KEY)
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

func SyncOrgsWithDB(orgs []*model.Org) error {

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

func FetchRepoOfOrgFromGitea(page int, orgName string, repos *[]*model.Repo) error {

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

	FetchRepoOfOrgFromGitea(page+1, orgName, repos)

	return nil
}

func SyncReposWithDB(repos []*model.Repo) error {
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

func FetchUsersFromGitea(page int, users *[]*model.User) error {
	currentPageUsers := []model.User{}

	url := fmt.Sprintf("https://gitea.vivasoftltd.com/api/v1/admin/users?page=%d&access_token=%s", page, config.AppConfig.GITEA.API_KEY)

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
		user.Last_updated = currentTime
		*users = append(*users, &user)
	}
	FetchUsersFromGitea(page+1, users)

	return nil
}
func SyncUsersWithDB(users []*model.User) error {
	collection := db.MongoDatabase.Collection(userCollection)

	uniqueUsers := make(map[int]*model.User)
	for _, user := range users {
		if _, exists := uniqueUsers[user.Id]; !exists {
			uniqueUsers[user.Id] = user
		}
	}
	for _, user := range uniqueUsers {
		fmt.Println(user.Username)
	}
	fmt.Println(len(uniqueUsers))

	documents := make([]interface{}, 0, len(uniqueUsers))

	for _, user := range uniqueUsers {
		documents = append(documents, user)
	}
	fmt.Println(len(documents))

	_, err := collection.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	return nil
}

func FetchUserActivityFromGitea(page int, userName string, activities *[]*model.Activity) error {
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
	FetchUserActivityFromGitea(page+1, userName, activities)

	return nil
}

func SyncActivitiesWithDB(activities []*model.Activity) error {
	collection := db.MongoDatabase.Collection(activitesCollection)

	documents := make([]interface{}, len(activities))
	repos := make(map[string]model.Repo)
	for i, user := range activities {
		repos[user.Repo.Name] = user.Repo
		documents[i] = user
	}

	if len(documents) == 0 {
		return nil
	}

	_, err := collection.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	username := activities[0].PerformedBy.Username
	collection = db.MongoDatabase.Collection(userCollection)
	filter := bson.M{"username": username}
	repoList := make([]model.Repo, 0, len(repos))
	for _, repo := range repos {
		repoList = append(repoList, repo)
	}

	update := bson.M{
		"$set": bson.M{
			"last_updated": time.Now(),
			"TotalCommits": len(activities),
			"repos":        repoList,
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func FetchDailyUserActivityFromGitea(page int, userName string, date string, activities *[]*model.Activity) error {
	currentPageActivities := []model.Activity{}
	limit := 50
	url := fmt.Sprintf("https://gitea.vivasoftltd.com/api/v1/users/%s/activities/feeds?only-performed-by=true&page=%d&limit=%d&date=%s&access_token=%s", userName, page, limit, date, config.AppConfig.GITEA.API_KEY)
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
	FetchDailyUserActivityFromGitea(page+1, userName, date, activities)

	return nil
}

func SyncDailyActivitiesWithDB(activities []*model.Activity) error {
	collection := db.MongoDatabase.Collection(activitesCollection)

	documents := make([]interface{}, len(activities))
	for i, user := range activities {
		documents[i] = user
	}

	if len(documents) == 0 {
		return nil
	}

	_, err := collection.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	username := activities[0].PerformedBy.Username
	collection = db.MongoDatabase.Collection(userCollection)
	filter := bson.M{"username": username}

	update := bson.M{
		"$set": bson.M{
			"last_updated": time.Now(),
			"TotalCommits": bson.M{"$sum": len(activities)},
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func ClearOrgs() error {
	collection := db.MongoDatabase.Collection(orgCollection)
	err := collection.Drop(context.Background())
	if err != nil {
		return err
	}
	return nil
}
func ClearRepos() error {
	collection := db.MongoDatabase.Collection(repoCollection)
	err := collection.Drop(context.Background())
	if err != nil {
		return err
	}
	return nil
}
func ClearUsers() error {
	collection := db.MongoDatabase.Collection(userCollection)
	err := collection.Drop(context.Background())
	if err != nil {
		return err
	}
	return nil
}
func ClearActivities() error {
	collection := db.MongoDatabase.Collection(activitesCollection)
	err := collection.Drop(context.Background())
	if err != nil {
		return err
	}
	return nil
}
