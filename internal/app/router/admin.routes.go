package router

import (
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/controller"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/middleware"
	"github.com/labstack/echo/v4"
)

// InitAuthRoutes sets up the routes for authentication
func InitAdminRoutes(e *echo.Echo) {
	adminGroup := e.Group("/admin")
	adminGroup.Use(middleware.Logger)
	e.GET("/admin/orgs", controller.GetAllOrganizations)
}
