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
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	client = &http.Client{}
)

const heatMapCollection = "heatmap"

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

func FetchUserActivityFromGitea(page int, userName string) ([]model.Activity, error) {
	activities := make([]model.Activity, 0)

	url := fmt.Sprintf("%s/users/%s/activities/feeds?only-performed-by=true&page=%d&access_token=%s", config.AppConfig.GITEA.Base_URL, userName, page, config.AppConfig.GITEA.API_KEY)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return activities, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return activities, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return activities, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&activities); err != nil {
		return activities, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if len(activities) == 0 {
		return activities, nil
	}
	commitActivities := make([]model.Activity, 0)
	for _, activity := range activities {
		if activity.OpType == "commit_repo" {
			commitActivities = append(commitActivities, activity)
		}
	}

	next_activities, err := FetchUserActivityFromGitea(page+1, userName)
	if err != nil {
		return activities, err
	}
	commitActivities = append(commitActivities, next_activities...)
	return commitActivities, nil
}

func SyncActivitiesWithDB(username string, activities []model.Activity) error {
	if len(activities) == 0 {
		collection := db.MongoDatabase.Collection(userCollection)
		filter := bson.M{"username": username}

		update := bson.M{
			"$set": bson.M{
				"last_updated": time.Now().UTC(),
			},
		}

		_, err := collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return err
		}
		return nil
	}
	collection := db.MongoDatabase.Collection(activitesCollection)

	documents := make([]interface{}, len(activities))
	repos := make(map[string]model.Repo)
	for i, user := range activities {
		repos[user.Repo.Name] = user.Repo
		documents[i] = user
	}
	_, err := collection.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	collection = db.MongoDatabase.Collection(userCollection)
	filter := bson.M{"username": username}
	repoList := make([]model.Repo, 0, len(repos))
	for _, repo := range repos {
		repoList = append(repoList, repo)
	}

	lastUpdateTime := time.Now().UTC()

	update := bson.M{
		"$set": bson.M{
			"last_updated":  lastUpdateTime,
			"total_commits": len(activities),
			"repos":         repoList,
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func FetchNewUserActivityFromGitea(page int, userName string, lastUpdateTime time.Time) ([]model.Activity, error) {

	activities := make([]model.Activity, 0)
	url := fmt.Sprintf("%s/users/%s/activities/feeds?only-performed-by=true&page=%d&access_token=%s", config.AppConfig.GITEA.Base_URL, userName, page, config.AppConfig.GITEA.API_KEY)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return activities, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return activities, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return activities, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&activities); err != nil {
		return activities, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	commitActivities := make([]model.Activity, 0)
	for _, activity := range activities {
		if activity.OpType == "commit_repo" && activity.Date.After(lastUpdateTime) {
			fmt.Println(userName, activity.Date, lastUpdateTime)
			commitActivities = append(commitActivities, activity)
		}
		if activity.Date.Before(lastUpdateTime) {
			return commitActivities, nil
		}
	}

	if len(commitActivities) == 0 {
		return commitActivities, nil
	}
	next_activities, err := FetchNewUserActivityFromGitea(page+1, userName, lastUpdateTime)
	if err != nil {
		return activities, err
	}
	commitActivities = append(commitActivities, next_activities...)
	return commitActivities, nil
}

func SyncNewActivitiesWithDB(username string, activities []model.Activity) error {
	if len(activities) == 0 {
		collection := db.MongoDatabase.Collection(userCollection)
		filter := bson.M{"username": username}

		update := bson.M{
			"$set": bson.M{
				"last_updated": time.Now().UTC(),
			},
		}

		_, err := collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return err
		}
		return nil
	}

	collection := db.MongoDatabase.Collection(activitesCollection)
	documents := make([]interface{}, len(activities))
	reposSet := make(map[string]model.Repo)
	for i, activity := range activities {
		documents[i] = activity
		reposSet[activity.Repo.Name] = activity.Repo
	}
	_, err := collection.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	collection = db.MongoDatabase.Collection(userCollection)
	filter := bson.M{"username": username}

	var user model.User
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return err
	}

	for _, repo := range user.Repos {
		reposSet[repo.Name] = repo
	}

	repoList := make([]model.Repo, 0, len(reposSet))
	for _, repo := range reposSet {
		repoList = append(repoList, repo)
	}
	fmt.Println("Total repos : ", len(reposSet), " for user ", username)

	update := bson.M{
		"$set": bson.M{
			"last_updated": time.Now().UTC(),
			"repos":        repoList,
		},
		"$inc": bson.M{
			"total_commits": len(activities),
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

func FetchUserHeatmapActivityFromGitea(userName string) (model.Heatmap, error) {
	heatMap := model.Heatmap{
		Username:       userName,
		HeatmapEntries: make([]model.HeatmapEntry, 0),
	}
	heatMapEntries := make([]model.HeatmapEntry, 0)
	heatmapURL := fmt.Sprintf("%s/users/%s/heatmap?access_token=%s", config.AppConfig.GITEA.Base_URL, userName, config.AppConfig.GITEA.API_KEY)

	heatmapReq, err := http.NewRequest("GET", heatmapURL, nil)
	if err != nil {
		return model.Heatmap{}, fmt.Errorf("failed to create HTTP request for heatmap: %w", err)
	}

	heatmapResp, err := client.Do(heatmapReq)
	if err != nil {
		return model.Heatmap{}, fmt.Errorf("failed to make HTTP request for heatmap: %w", err)
	}
	defer heatmapResp.Body.Close()

	if heatmapResp.StatusCode != http.StatusOK {
		return model.Heatmap{}, fmt.Errorf("unexpected status code for heatmap: %d", heatmapResp.StatusCode)
	}

	if err := json.NewDecoder(heatmapResp.Body).Decode(&heatMapEntries); err != nil {
		return model.Heatmap{}, fmt.Errorf("failed to decode JSON response for heatmap: %w", err)
	}
	heatMap.HeatmapEntries = heatMapEntries
	return heatMap, nil
}

func SyncHeatMaps(heatmap model.Heatmap) error {
	collection := db.MongoDatabase.Collection(heatMapCollection)

	_, err := collection.InsertOne(context.Background(), heatmap)
	if err != nil {
		return fmt.Errorf("failed to insert heatmap: %w", err)
	}

	return nil
}

func ClearHeatmaps() error {
	collection := db.MongoDatabase.Collection(heatMapCollection)
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
