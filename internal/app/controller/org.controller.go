package controller

import (
	"net/http"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"

	"github.com/labstack/echo/v4"
)

func GetAllOrganizationsFromExternalAPI(c echo.Context) error {
	var orgs []*model.Org

	err := repository.GetAllOrganizationsFromGitea(1, &orgs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, orgs)
}

func GetAllRepoOfOrganization(c echo.Context) error {
	orgName := c.Param("org")
	var repos []*model.Repo

	err := repository.GetAllRepoOfOrganization(1, orgName, &repos)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, repos)
}

func GetAllRepoFromDB(c echo.Context) error {
	var repos []*model.Repo
	org := c.Param("org")

	err := repository.GetAllRepoFromDB(org, &repos)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, repos)
}
