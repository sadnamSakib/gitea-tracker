package router

import (
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/controller"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/middleware"
	"github.com/labstack/echo/v4"
)

// InitAuthRoutes sets up the routes for authentication
func InitUserRoutes(e *echo.Echo) {
	userGroup := e.Group("/user")
	userGroup.Use(middleware.Logger)
	e.GET("/users/:org", controller.GetAllUsers)

}
