package router

import (
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/controller"
	"github.com/labstack/echo/v4"
)

// InitAuthRoutes sets up the routes for authentication
func InitOrgRoutes(e *echo.Echo) {
	orgGroup := e.Group("/org")
	orgGroup.GET("/orgs", controller.GetAllOrganizationsFromExternalAPI)
	orgGroup.GET("/orgs/:org/repos", controller.GetAllRepoFromDB)

}
