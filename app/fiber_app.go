package app

import (
	"log"
	"openidea-idea-social-media-app/customErr"
	"openidea-idea-social-media-app/db"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func StartFiberApp(port string, prefork bool) {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		Prefork:     prefork,
	})

	prometheus := fiberprometheus.New("social-app-service")
	prometheus.RegisterAt(app, "/metrics")

	dbPool := db.GetConnectionPool()
	defer dbPool.Close()

	app.Use(prometheus.Middleware)
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	RegisterRoute(app, dbPool)

	app.Use(customErr.NotFoundHandler)

	err := app.Listen(":" + port)
	log.Fatal(err)
}
