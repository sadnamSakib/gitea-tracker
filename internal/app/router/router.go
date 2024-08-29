package router

import (
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	InitUserRoutes(e)
	InitAdminRoutes(e)
}
