package controller

import (
	"net/http"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/service"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"

	"github.com/labstack/echo/v4"
)

func GetUser(c echo.Context) error {
	userName := c.Param("username")
	user, err := repository.GetUser(userName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}
func FollowUser(c echo.Context) error {
	userName := c.Param("username")
	err := repository.FollowUser(userName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "User Followed")
}
func UnfollowUser(c echo.Context) error {
	userName := c.Param("username")
	err := repository.UnfollowUser(userName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "User Unfollowed")
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

func GetAllUsers(c echo.Context) error {
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	users, err := repository.GetAllUsers(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

func SearchUsers(c echo.Context) error {
	query := c.QueryParam("query")
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	users, err := repository.SearchUsers(query, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)

}

func SyncUsers(c echo.Context) error {
	usersSynced, err := service.SyncAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"Users": usersSynced})
}
