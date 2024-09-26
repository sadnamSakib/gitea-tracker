package main

import (
	"log"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/router"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/service"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/config"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	config.LoadConfig()

	db.Connect()
	defer db.Disconnect()

	e := echo.New()
	e.Static("/web", "web")

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} [${method}] ${uri} (${status}) (${latency_human})\n",
	}))
	e.Use(middleware.Recover())

	router.InitRoutes(e)

	c := service.InitCronScheduler()
	defer c.Stop()

	log.Fatal(e.Start(":8080"))
}
