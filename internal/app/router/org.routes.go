package router

import (
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/controller"
	"github.com/labstack/echo/v4"
)

// InitAuthRoutes sets up the routes for authentication
func InitOrgRoutes(e *echo.Group) {
	orgGroup := e.Group("/orgs")
	orgGroup.GET("/", controller.GetAllOrganizationFromDB)
	orgGroup.GET("/:org/repos/search", controller.SearchRepos)
	orgGroup.GET("/:org/repos", controller.GetAllRepoFromDB)
	orgGroup.GET("/:org/repos/:repo/users", controller.GetAllUsersOfRepo)
	orgGroup.GET("/:org/repos/:repo/users/search", controller.SearchUsersOfRepo)
	orgGroup.POST("/:org/repos/:repo/follow", controller.FollowRepo)
	orgGroup.POST("/:org/repos/:repo/unfollow", controller.UnfollowRepo)

}
