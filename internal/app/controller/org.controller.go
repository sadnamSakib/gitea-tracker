package controller

import (
	"net/http"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"

	"github.com/labstack/echo/v4"
)

func GetAllOrganizationsFromExternalAPI(c echo.Context) error {
	var orgs []*model.Org

	err := repository.FetchOrgsFromGitea(1, &orgs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, orgs)
}

func GetAllOrganizationFromDB(c echo.Context) error {
	var orgs []*model.Org

	err := repository.GetAllOrgs(&orgs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, orgs)
}

func GetAllRepoOfOrganization(c echo.Context) error {
	orgName := c.Param("org")
	var repos []*model.Repo

	err := repository.GetAllReposFromOrg(orgName, &repos)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, repos)
}

func GetAllRepoFromDB(c echo.Context) error {
	var repos []*model.Repo
	org := c.Param("org")

	err := repository.GetAllReposFromOrg(org, &repos)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, repos)
}

func GetAllUsersOfRepo(c echo.Context) error {
	org := c.Param("org")
	repo := c.Param("repo")
	var users []*model.User
	err := repository.GetAllUsersFromRepo(org, repo, &users)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}
