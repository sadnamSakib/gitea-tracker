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

func SyncAllUsers() error {
	users, err := repository.FetchUsersFromGitea(1)
	if err != nil {
		return err
	}
	err = repository.ClearUsers()
	if err != nil {
		return err
	}
	err = repository.SyncUsersWithDB(users)
	if err != nil {
		return err
	}

	return nil
}

func SyncUsers(c echo.Context) error {
	err := SyncAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Users Synced")
}

func SyncUserHeatmap(username string) error {
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

func SyncAllHeatmaps() error {
	users, err := repository.GetAllUsers("", "")
	if err != nil {
		return err
	}
	err = repository.ClearHeatmaps()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 10)
	errorsChan := make(chan error, len(users))
	userSynced := 0
	for _, user := range users {
		wg.Add(1)
		sem <- struct{}{}
		go func(username string) {
			defer wg.Done()
			defer func() { <-sem }()

			if err := SyncUserHeatmap(username); err != nil {
				errorsChan <- err
			} else {
				userSynced++
			}
		}(user.Username)
	}
	wg.Wait()
	close(errorsChan)
	for e := range errorsChan {
		err = e
		fmt.Println(err)
	}
	fmt.Printf("Synced %d users\n", userSynced)
	return err
}

func SyncHeatMaps(c echo.Context) error {
	err := SyncAllHeatmaps()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Heatmaps Synced")
}

func SyncUserActivities(userName string) error {

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

func SyncAllActivities() error {

	users, err := repository.GetAllUsers("", "")
	if err != nil {
		return err
	}
	err = repository.ClearActivities()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 10)
	errorsChan := make(chan error, len(users))
	userSynced := 0
	for _, user := range users {
		wg.Add(1)
		sem <- struct{}{}
		go func(username string) {
			defer wg.Done()
			defer func() { <-sem }()

			if err := SyncUserActivities(username); err != nil {
				errorsChan <- err
			} else {
				userSynced++
			}
		}(user.Username)
	}
	wg.Wait()
	close(errorsChan)
	for e := range errorsChan {
		err = e
		fmt.Println(err)
	}
	fmt.Printf("Synced %d users\n", userSynced)
	return err
}

func SyncActivities(c echo.Context) error {
	err := SyncAllActivities()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Activities Synced")
}

func SyncAllOrganizations() error {
	orgs, err := repository.FetchOrgsFromGitea(1)
	if err != nil {
		return err
	}
	err = repository.ClearOrgs()
	if err != nil {
		return err
	}
	err = repository.SyncOrgsWithDB(orgs)
	if err != nil {
		return err
	}

	return nil
}

func SyncOrganizations(c echo.Context) error {
	err := SyncAllOrganizations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Organizations Synced")
}

func SyncOrgRepos(orgName string, wg *sync.WaitGroup) error {
	defer wg.Done()

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

func SyncAllRepos() error {
	var orgs []*model.Org
	orgs, err := repository.GetAllOrgs()
	if err != nil {

		return err
	}
	err = repository.ClearRepos()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	for _, org := range orgs {
		wg.Add(1)
		go SyncOrgRepos(org.Username, &wg)
	}
	wg.Wait()
	return nil
}
func SyncRepos(c echo.Context) error {
	err := SyncAllRepos()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Repos Synced")
}

func SyncNewUserActivities(username string) error {

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

func SyncDailyActivities() error {

	users, err := repository.GetAllUsers("", "")
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 10)
	errorsChan := make(chan error, len(users))
	usersSynced := 0
	for _, user := range users {
		wg.Add(1)
		sem <- struct{}{}
		go func(username string) {
			defer wg.Done()
			defer func() { <-sem }()
			if err := SyncNewUserActivities(username); err != nil {
				errorsChan <- err
			} else {
				usersSynced++
			}
		}(user.Username)
	}
	wg.Wait()
	close(errorsChan)
	for e := range errorsChan {
		err = e
		fmt.Println(err)
	}
	fmt.Println("Synced ", usersSynced, " users")
	return err
}

func SyncNewActivity(c echo.Context) error {
	err := SyncDailyActivities()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "New Activities Synced")
}

func DailySync(c echo.Context) error {
	err := SyncAllOrganizations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("Organizations Synchronised")

	err = SyncAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("Users Synchronised")

	err = SyncAllRepos()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("Repos Synchronised")
	err = SyncAllHeatmaps()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = SyncDailyActivities()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("Daily Activities Synchronised")
	return c.JSON(http.StatusOK, "Daily Sync Completed")
}
