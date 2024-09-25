package router

import (
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/middleware"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	e.Use(middleware.Logger)
	apiGroup := e.Group("/api")
	InitOrgRoutes(apiGroup)
	InitUserRoutes(apiGroup)
	InitSyncRoutes(apiGroup)
	InitWebRoutes(e)
}
