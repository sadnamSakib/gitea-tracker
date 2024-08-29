package main

import (
	"log"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/router"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/config"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load application configuration
	config.LoadConfig()

	// Initialize MongoDB connection
	db.Connect()
	defer db.Disconnect()

	// Create a new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} [${method}] ${uri} (${status})\n",
	}))
	e.Use(middleware.Recover())

	// main.go or wherever you initialize Echo
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"}, // Adjust to match your frontend
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowCredentials: true, // This is necessary for cookies to work
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Initialize API routes
	router.InitRoutes(e)

	log.Fatal(e.Start(":8080"))
}
