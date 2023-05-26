package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/saqura-io/torii/config"
	"github.com/saqura-io/torii/internal/router"
	"log"
)

func main() {
	app := fiber.New()
	app.Use(recover2.New())
	app.Use(logger.New())

	routeConfig, err := router.LoadConfig("./config/services.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	router.SetupRoutes(app, routeConfig)

	port := config.Get("SERVER_PORT", "8080")
	err = app.Listen(":" + port)

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
