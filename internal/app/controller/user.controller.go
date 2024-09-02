package controller

import (
	"net/http"

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
