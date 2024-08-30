package controller

import (
	"fmt"
	"net/http"
	"sync"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"

	"github.com/labstack/echo/v4"
)

func GetAllUsersFromExternalAPI(c echo.Context) error {
	var users []*model.User
	var orgName string = c.Param("org")

	err := repository.GetAllUserFromGitea(1, orgName, &users)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func GetUserActivityFromExternalAPI(userName string) ([]*model.Activity, error) {
	var activities []*model.Activity

	err := repository.GetAllUserActivityFromGitea(1, userName, &activities)

	if err != nil {
		return nil, err
	}
	return activities, nil
}

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

func SyncUserActivities(userName string, wg *sync.WaitGroup, sem chan struct{}) error {
	defer wg.Done()
	defer func() { <-sem }()

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
	sem := make(chan struct{}, len(users))
	for _, user := range users {
		wg.Add(1)
		sem <- struct{}{}
		go SyncUserActivities(user.Username, &wg, sem)
	}
	wg.Wait()
	return nil
}

func SyncActivities(c echo.Context) error {
	err := SyncAllActivities()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Activities Synced")
}
