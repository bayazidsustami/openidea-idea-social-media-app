package app

import (
	"log"
	"openidea-idea-social-media-app/config"
	"openidea-idea-social-media-app/customErr"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func StartFiberApp() {
	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		IdleTimeout:  config.IdleTimeout,
		WriteTimeout: config.WriteTimeout,
		ReadTimeout:  config.ReadTimeout,
		Prefork:      true,
		ErrorHandler: customErr.ErrorHandler,
	})

	app.Use(logger.New())

	RegisterRoute(app)

	err := app.Listen("localhost:8000")
	log.Fatal(err)
}
