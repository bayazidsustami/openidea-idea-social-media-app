package app

import (
	"log"
	"openidea-idea-social-media-app/config"
	"openidea-idea-social-media-app/customErr"
	"openidea-idea-social-media-app/db"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func StartFiberApp() {
	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		IdleTimeout:  config.IdleTimeout,
		WriteTimeout: config.WriteTimeout,
		ReadTimeout:  config.ReadTimeout,
		Prefork:      true,
	})

	dbPool := db.GetConnectionPool()
	defer dbPool.Close()

	app.Use(logger.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	RegisterRoute(app, dbPool)

	app.Use(customErr.NotFoundHandler)

	err := app.Listen("localhost:8000")
	log.Fatal(err)
}
