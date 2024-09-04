package controller

import (
	"net/http"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"

	"github.com/labstack/echo/v4"
)

func GetAllOrganizationFromDB(c echo.Context) error {

	orgs, err := repository.GetAllOrgs()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, orgs)
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
