package router

import (
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/controller"
	"github.com/labstack/echo/v4"
)

// InitWebRoutes sets up the routes for web
func InitWebRoutes(e *echo.Echo) {
	webGroup := e.Group("")
	webGroup.GET("", controller.RenderHome)
	webGroup.GET("/orgs", controller.RenderOrganizations)
	webGroup.GET("/orgs/:org/repos", controller.RenderRepos)
	webGroup.GET("/users", controller.RenderUsers)
	webGroup.GET("/users/:user", controller.RenderUser)

}
