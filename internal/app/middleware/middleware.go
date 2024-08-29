package middleware

import (
	"github.com/labstack/echo/v4"
)

// AuthMiddleware - Example middleware function for authentication
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Perform some action
		c.Logger().Info("Request received")
		return next(c)
	}
}
