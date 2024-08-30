package router

import (
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/controller"
	"github.com/labstack/echo/v4"
)

// InitAuthRoutes sets up the routes for authentication
func InitAdminRoutes(e *echo.Echo) {
	adminGroup := e.Group("/admin")
	adminGroup.GET("/orgs", controller.GetAllOrganizationsFromExternalAPI)
}
