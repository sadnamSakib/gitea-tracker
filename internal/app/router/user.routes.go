package router

import (
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/controller"
	"github.com/labstack/echo/v4"
)

// InitAuthRoutes sets up the routes for authentication
func InitUserRoutes(e *echo.Group) {
	userGroup := e.Group("/users")
	userGroup.GET("/", controller.GetAllUsers)
	userGroup.GET("/search", controller.SearchUsers)
	userGroup.GET("/:username", controller.GetUser)
	userGroup.GET("/:username/activities", controller.GetUserActivityByDateRange)
	userGroup.POST("/:username/follow", controller.FollowUser)
	userGroup.POST("/:username/unfollow", controller.UnfollowUser)
}
