package controller

import (
	"net/http"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/service"
	"github.com/labstack/echo/v4"
)

func SyncNewActivity(c echo.Context) error {
	users, err := service.SyncDailyActivities()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"User's New Activities": users})
}
