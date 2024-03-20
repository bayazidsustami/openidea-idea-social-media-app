package app

import (
	"openidea-idea-social-media-app/controller"
	"openidea-idea-social-media-app/repository"
	"openidea-idea-social-media-app/security"
	"openidea-idea-social-media-app/service"
	"openidea-idea-social-media-app/validation"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoute(app *fiber.App, dbPool *pgxpool.Pool) {

	validator := validator.New()
	validation.RegisterValidation(validator)

	authService := service.NewAuthService()

	userRepository := repository.New()
	userService := service.NewUserService(userRepository, validator, dbPool, authService)
	userController := controller.New(userService)

	friendRepository := repository.NewFriendRepository(dbPool)
	friendService := service.NewFriendsService(validator, friendRepository)
	friendController := controller.NewFriendsController(friendService, authService)

	user := app.Group("/v1/user")
	user.Post("/register", userController.Register)
	user.Post("/login", userController.Login)

	app.Use(security.CheckTokenHeaderExist)
	app.Use(security.GetJwtTokenHandler())

	app.Post("v1/friend", friendController.AddFriend)
	app.Delete("v1/friend", friendController.RemoveFriends)
}
