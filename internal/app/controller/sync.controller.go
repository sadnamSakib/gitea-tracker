package controller

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
	"github.com/labstack/echo/v4"
)

func SyncAllUsers() (int, error) {
	users, err := repository.FetchUsersFromGitea(1)
	if err != nil {
		return 0, err
	}

	err = repository.SyncUsersWithDB(users)
	if err != nil {
		return 0, err
	}

	return len(users), nil
}

func SyncUsers(c echo.Context) error {
	usersSynced, err := SyncAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"Users": usersSynced})
}

func SyncUserHeatmap(username string, wg *sync.WaitGroup, sem chan struct{}) error {
	defer wg.Done()
	defer func() { <-sem }()
	heatmap, err := repository.FetchUserHeatmapActivityFromGitea(username)
	if err != nil {
		return err
	}
	err = repository.SyncHeatMaps(heatmap)
	if err != nil {
		return err
	}
	fmt.Printf("User %s's Heatmap synced \n", username)
	return nil
}

func SyncAllHeatmaps() (int, error) {
	users, err := repository.GetAllUsers("", "")
	if err != nil {
		return 0, err
	}
	err = repository.ClearHeatmaps()
	if err != nil {
		return 0, err
	}
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 10)
	errorsChan := make(chan error, len(users))
	userSynced := 0
	for _, user := range users {
		wg.Add(1)
		sem <- struct{}{}

		if err := SyncUserHeatmap(user.Username, &wg, sem); err != nil {
			errorsChan <- err
		} else {
			userSynced++
		}

	}
	wg.Wait()
	close(errorsChan)
	close(sem)
	for e := range errorsChan {
		err = e
		fmt.Println(err)
	}
	fmt.Printf("Synced %d users\n", userSynced)
	return len(users), err
}

func SyncHeatMaps(c echo.Context) error {
	users, err := SyncAllHeatmaps()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"User's Heatmaps": users})
}

func SyncUserActivities(userName string, wg *sync.WaitGroup, sem chan struct{}) error {
	defer wg.Done()
	defer func() { <-sem }()
	activities, err := repository.FetchUserActivityFromGitea(1, userName)
	if err != nil {
		return err
	}
	err = repository.SyncActivitiesWithDB(userName, activities)
	if err != nil {
		return err
	}
	fmt.Printf("User %s synced with %d activities\n", userName, len(activities))
	return nil
}

func SyncAllActivities() (int, error) {

	users, err := repository.GetAllUsers("", "")
	if err != nil {
		return 0, err
	}
	err = repository.ClearActivities()
	if err != nil {
		return 0, err
	}

	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 10)
	errorsChan := make(chan error, len(users))
	userSynced := 0
	for _, user := range users {
		wg.Add(1)
		sem <- struct{}{}

		if err := SyncUserActivities(user.Username, &wg, sem); err != nil {
			errorsChan <- err
		} else {
			userSynced++
		}

	}
	wg.Wait()
	close(errorsChan)
	close(sem)
	for e := range errorsChan {
		err = e
		fmt.Println(err)
	}

	return len(users), err
}

func SyncActivities(c echo.Context) error {
	users, err := SyncAllActivities()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"Users Activities": users})
}

func SyncAllOrganizations() (int, error) {
	orgs, err := repository.FetchOrgsFromGitea(1)
	if err != nil {
		return 0, err
	}
	err = repository.ClearOrgs()
	if err != nil {
		return 0, err
	}
	err = repository.SyncOrgsWithDB(orgs)
	if err != nil {
		fmt.Println("Error here")
		return 0, err
	}
	fmt.Println(orgs, len(orgs))
	return len(orgs), nil
}

func SyncOrganizations(c echo.Context) error {
	orgsSynced, err := SyncAllOrganizations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"Organizations": orgsSynced})
}

func SyncOrgRepos(orgName string, wg *sync.WaitGroup, sem chan struct{}) error {
	defer wg.Done()
	defer func() { <-sem }()

	repos, err := repository.FetchRepoOfOrgFromGitea(1, orgName)
	if err != nil {
		return err
	}

	err = repository.SyncReposWithDB(repos)
	if err != nil {
		return err
	}
	fmt.Printf("Synced %d repositories of %s\n", len(repos), orgName)

	return nil
}

func SyncAllRepos() (int, error) {
	var orgs []*model.Org
	orgs, err := repository.GetAllOrgs()
	if err != nil {

		return 0, err
	}
	err = repository.ClearRepos()
	if err != nil {
		return 0, err
	}
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 10)
	errorsChan := make(chan error, len(orgs))
	for _, org := range orgs {
		wg.Add(1)
		sem <- struct{}{}
		if err := SyncOrgRepos(org.Username, &wg, sem); err != nil {
			errorsChan <- err
		}
	}
	close(errorsChan)
	close(sem)
	for e := range errorsChan {
		err = e
		fmt.Println(err)
	}
	wg.Wait()
	return len(orgs), nil
}
func SyncRepos(c echo.Context) error {
	orgsSynced, err := SyncAllRepos()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"Repos For Organization": orgsSynced})
}

func SyncNewUserActivities(username string, wg *sync.WaitGroup, sem chan struct{}) error {
	defer wg.Done()
	defer func() { <-sem }()
	format := "2006-01-02"
	currentDate := time.Now().Format(format)
	user, err := repository.GetUser(username)

	if err != nil {
		return err
	}
	lastUpdateTime := user.Last_updated
	activities, err := repository.FetchNewUserActivityFromGitea(1, username, currentDate, lastUpdateTime)
	if err != nil {
		return err
	}
	err = repository.SyncNewActivitiesWithDB(username, activities)
	if err != nil {
		return err
	}
	fmt.Printf("User %s has %d new activities.\n", username, len(activities))
	return nil
}

func SyncDailyActivities() (int, error) {

	users, err := repository.GetAllUsers("", "")
	if err != nil {
		return 0, err
	}

	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 10)
	errorsChan := make(chan error, len(users))
	usersSynced := 0
	for _, user := range users {
		wg.Add(1)
		sem <- struct{}{}

		if err := SyncNewUserActivities(user.Username, &wg, sem); err != nil {
			errorsChan <- err
		} else {
			usersSynced++
		}

	}
	wg.Wait()
	close(errorsChan)
	close(sem)
	for e := range errorsChan {
		err = e
		fmt.Println(err)
	}
	fmt.Println("Synced ", usersSynced, " users")
	return len(users), err
}

func SyncNewActivity(c echo.Context) error {
	users, err := SyncDailyActivities()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"User's New Activities": users})
}

func DailySync(c echo.Context) error {
	orgs, err := SyncAllOrganizations()

	if err != nil {

		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("Organizations Synchronised")

	users, err := SyncAllUsers()
	if err != nil {

		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("Users Synchronised")

	orgsRepo, err := SyncAllRepos()
	if err != nil {

		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("Repos Synchronised")

	usersHeatmap, err := SyncAllHeatmaps()
	if err != nil {

		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	usersNewActivities, err := SyncDailyActivities()
	if err != nil {

		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("Daily Activities Synchronised")
	return c.JSON(http.StatusOK, map[string]int{"Organizations Synced": orgs, "Users Synced": users, "Repos Synced For Organizations": orgsRepo, "Heatmaps Synced For Users": usersHeatmap, "New Activities Synced For Users": usersNewActivities})
}
