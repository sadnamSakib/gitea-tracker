package controller

import (
	"net/http"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/service"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"
	"github.com/labstack/echo/v4"
)

func SyncRepos(c echo.Context) error {
	orgsSynced, err := service.SyncAllRepos()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"Repos For Organization": orgsSynced})
}
func GetAllRepoFromDB(c echo.Context) error {

	org := c.Param("org")
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")

	repos, err := repository.GetAllReposFromOrg(org, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, repos)
}

func GetAllUsersOfRepo(c echo.Context) error {
	org := c.Param("org")
	repo := c.Param("repo")
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")

	users, err := repository.GetAllUsersFromRepo(org, repo, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}
func FollowRepo(c echo.Context) error {
	org := c.Param("org")
	repo := c.Param("repo")
	err := repository.FollowRepo(org, repo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Repo Followed")
}
func UnfollowRepo(c echo.Context) error {
	org := c.Param("org")
	repo := c.Param("repo")
	err := repository.UnfollowRepo(org, repo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Repo Unfollowed")
}

func SearchRepos(c echo.Context) error {
	org := c.Param("org")
	query := c.QueryParam("query")
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")

	repos, err := repository.SearchRepos(org, query, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, repos)
}

func SearchUsersOfRepo(c echo.Context) error {
	org := c.Param("org")
	repo := c.Param("repo")
	query := c.QueryParam("query")
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")

	users, err := repository.SearchUsersOfRepo(org, repo, query, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}
