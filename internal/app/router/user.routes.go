package router

import (
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/controller"
	"github.com/labstack/echo/v4"
)

// InitAuthRoutes sets up the routes for authentication
func InitUserRoutes(e *echo.Echo) {
	userGroup := e.Group("/users")
	userGroup.GET("/:org", controller.GetAllUsersFromExternalAPI)

}
