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

func GetUser(c echo.Context) error {
	user := &model.User{}
	userName := c.Param("username")
	err := repository.GetUser(userName, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

func GetUserActivityByDateRange(c echo.Context) error {
	userName := c.Param("username")
	start_date := c.QueryParam("start_date")
	end_date := c.QueryParam("end_date")
	count_only := c.QueryParam("count_only")
	repo := c.QueryParam("repo")
	activities, err := repository.GetUserActivityByDateRange(userName, start_date, end_date, repo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if count_only == "true" {
		return c.JSON(http.StatusOK, len(activities))
	} else {
		return c.JSON(http.StatusOK, activities)
	}
}
