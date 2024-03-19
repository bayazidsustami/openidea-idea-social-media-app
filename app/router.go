package app

import (
	"openidea-idea-social-media-app/controller"
	"openidea-idea-social-media-app/db"
	"openidea-idea-social-media-app/repository"
	"openidea-idea-social-media-app/security"
	"openidea-idea-social-media-app/service"
	"openidea-idea-social-media-app/validation"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoute(app *fiber.App) {

	validator := validator.New()
	validation.RegisterValidation(validator)

	dbPool := db.GetConnectionPool()

	authService := service.NewAuthService()

	userRepository := repository.New()
	userService := service.NewUserService(userRepository, validator, dbPool, authService)
	userController := controller.New(userService)

	api := app.Group("/v1")
	api.Post("/user/register", userController.Register)

	app.Use(security.CheckTokenHeaderExist)
	app.Use(security.GetJwtTokenHandler())

}
