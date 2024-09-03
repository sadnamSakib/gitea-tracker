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
	var users []*model.User
	err := repository.FetchUsersFromGitea(1, &users)
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
	var activities []*model.Activity
	err := repository.FetchUserActivityFromGitea(1, userName, &activities)
	if err != nil {
		return err
	}
	err = repository.SyncActivitiesWithDB(activities)
	if err != nil {
		return err
	}
	fmt.Println("userName: ", userName, " Synced")
	return nil
}

func SyncAllActivities() error {
	var users []*model.User
	err := repository.GetAllUsers(&users)
	if err != nil {
		return err
	}
	err = repository.ClearActivities()
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

			if err := SyncUserActivities(username); err != nil {
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

func SyncActivities(c echo.Context) error {

	err := SyncAllActivities()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Activities Synced")
}

func SyncOrganizations(c echo.Context) error {
	var orgs []*model.Org

	err := repository.FetchOrgsFromGitea(1, &orgs)
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
	var repos []*model.Repo

	err := repository.FetchRepoOfOrgFromGitea(1, orgName, &repos)
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
	err := repository.GetAllOrgs(&orgs)
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
	var activities []*model.Activity
	format := "2006-01-02"
	currentDate := time.Now().Format(format)
	err := repository.FetchDailyUserActivityFromGitea(1, username, currentDate, &activities)
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
	var users []*model.User
	err := repository.GetAllUsers(&users)
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
