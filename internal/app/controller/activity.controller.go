package controller

import (
	"net/http"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/service"
	"github.com/labstack/echo/v4"
)

func SyncActivities(c echo.Context) error {
	users, err := service.SyncAllActivities()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"Users Activities": users})
}
