package app

import (
	"log"
	"openidea-idea-social-media-app/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func StartFiberApp() {
	app := fiber.New(fiber.Config{
		IdleTimeout:  config.IdleTimeout,
		WriteTimeout: config.WriteTimeout,
		ReadTimeout:  config.ReadTimeout,
		Prefork:      true,
	})

	app.Use(logger.New())

	err := app.Listen("localhost:8000")
	log.Fatal(err)
}
