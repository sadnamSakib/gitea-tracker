package router

import (
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/controller"
	"github.com/labstack/echo/v4"
)

// InitAuthRoutes sets up the routes for authentication
func InitSyncRoutes(e *echo.Group) {
	syncGroup := e.Group("/sync")
	syncGroup.GET("/orgs", controller.SyncOrganizations)
	syncGroup.GET("/users", controller.SyncUsers)
	syncGroup.GET("/activities", controller.SyncActivities)
	syncGroup.GET("/repos", controller.SyncRepos)
	syncGroup.GET("/newActivity", controller.SyncNewActivity)
	syncGroup.GET("/dailySync", controller.DailySync)
	syncGroup.GET("/totalSync", controller.TotalSync)
}
