package controller

import (
	"net/http"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/service"
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

func SyncOrganizations(c echo.Context) error {
	orgsSynced, err := service.SyncAllOrganizations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"Organizations": orgsSynced})
}
