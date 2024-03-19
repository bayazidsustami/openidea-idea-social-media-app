package controller

import (
	"log"
	"openidea-idea-social-media-app/customErr"
	user_model "openidea-idea-social-media-app/models/user"
	"openidea-idea-social-media-app/service"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserService service.UserService
}

func New(
	userService service.UserService,
) UserController {
	return UserController{
		UserService: userService,
	}
}

func (controller *UserController) Register(ctx *fiber.Ctx) error {
	userRequest := new(user_model.UserRegisterRequest)

	err := ctx.BodyParser(userRequest)
	if err != nil {
		log.Fatal(customErr.ErrorBadRequest)
	}

	response := controller.UserService.Register(ctx.UserContext(), *userRequest)

	return ctx.Status(fiber.StatusCreated).JSON(response)
}
