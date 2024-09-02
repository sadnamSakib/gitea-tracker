package controller

import (
	"fmt"
	"net/http"
	"sync"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
	"github.com/labstack/echo/v4"
)

func SyncOrgUsers(orgName string, wg *sync.WaitGroup) error {
	defer wg.Done()
	var users []*model.User

	err := repository.GetAllUserFromGitea(1, orgName, &users)
	if err != nil {
		return err
	}
	fmt.Println("users:", users)
	err = repository.SyncUsers(users)
	if err != nil {
		return err
	}

	return nil
}

func SyncAllOrgUsers() error {
	var orgs []*model.Org
	err := repository.GetAllOrganizationFromDB(&orgs)
	if err != nil {

		return err
	}
	fmt.Println("orgs:", orgs)
	wg := sync.WaitGroup{}
	for _, org := range orgs {
		wg.Add(1)
		go SyncOrgUsers(org.Username, &wg)
	}
	go func() { wg.Wait() }()
	return nil
}

func SyncUsers(c echo.Context) error {
	err := SyncAllOrgUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Users Synced")
}

func SyncUserActivities(userName string) error {

	activities, err := GetUserActivityFromExternalAPI(userName)
	if err != nil {
		return err
	}
	err = repository.SyncActivities(activities)
	if err != nil {
		return err
	}
	fmt.Println("userName: ", userName, " Synced")
	return nil
}

func SyncAllActivities() error {
	var users []*model.User
	err := repository.GetAllUsersFromDB(&users)
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

	err := repository.GetAllOrganizationsFromGitea(1, &orgs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	err = repository.SyncOrganizations(orgs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, orgs)
}
