package router

import (
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/controller"
	"github.com/labstack/echo/v4"
)

// InitAuthRoutes sets up the routes for authentication
func InitUserRoutes(e *echo.Echo) {
	userGroup := e.Group("/api/users")
	userGroup.GET("/", controller.GetAllUsers)
	userGroup.GET("/search", controller.SearchUsers)
	userGroup.GET("/:username", controller.GetUser)
	userGroup.GET("/:username/activities", controller.GetUserActivityByDateRange)
	userGroup.GET("/:username/heatmap", controller.GetUserHeatmap)
}
