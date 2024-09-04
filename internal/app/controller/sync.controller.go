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

func SyncUsers(c echo.Context) error {

	users, err := repository.FetchUsersFromGitea(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = repository.ClearUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = repository.SyncUsersWithDB(users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("%d users synced", len(users)))
}

func SyncUserActivities(userName string) error {

	activities, err := repository.FetchUserActivityFromGitea(1, userName)
	if err != nil {
		return err
	}
	err = repository.SyncActivitiesWithDB(activities)
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

func SyncOrganizations(c echo.Context) error {
	orgs, err := repository.FetchOrgsFromGitea(1)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = repository.ClearOrgs()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = repository.SyncOrgsWithDB(orgs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, orgs)
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

func SyncDailyUserActivities(username string) error {

	format := "2006-01-02"
	currentDate := time.Now().Format(format)
	activities, err := repository.FetchDailyUserActivityFromGitea(1, username, currentDate)
	if err != nil {
		return err
	}
	err = repository.SyncDailyActivitiesWithDB(activities)
	if err != nil {
		return err
	}
	fmt.Println("userName: ", username, " Synced")
	return nil
}

func SyncDailyActivities() error {

	users, err := repository.GetAllUsers("", "")
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 8)
	errorsChan := make(chan error, len(users))

	for _, user := range users {
		wg.Add(1)
		sem <- struct{}{}
		go func(username string) {
			defer wg.Done()
			defer func() { <-sem }()
			if err := SyncDailyUserActivities(username); err != nil {
				errorsChan <- err
			}
		}(user.Username)
	}
	wg.Wait()
	close(errorsChan)
	for e := range errorsChan {
		err = e
	}
	return err
}
