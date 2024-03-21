package controller

import (
	"openidea-idea-social-media-app/customErr"
	user_model "openidea-idea-social-media-app/models/user"
	"openidea-idea-social-media-app/service"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserService service.UserService
	AuthService service.AuthService
}

func New(
	userService service.UserService,
	authService service.AuthService,
) UserController {
	return UserController{
		UserService: userService,
		AuthService: authService,
	}
}

func (controller *UserController) Register(ctx *fiber.Ctx) error {
	userRequest := new(user_model.UserRegisterRequest)

	err := ctx.BodyParser(userRequest)
	if err != nil {
		return customErr.ErrorBadRequest
	}

	response, err := controller.UserService.Register(ctx.UserContext(), *userRequest)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (controller *UserController) Login(ctx *fiber.Ctx) error {
	userRequest := new(user_model.UserLoginRequest)

	err := ctx.BodyParser(userRequest)
	if err != nil {
		return customErr.ErrorBadRequest
	}

	response, err := controller.UserService.Login(ctx.UserContext(), *userRequest)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}

func (controller *UserController) UpdateEmail(ctx *fiber.Ctx) error {
	emailReq := new(user_model.UpdateEmailRequest)

	err := ctx.BodyParser(emailReq)
	if err != nil {
		return customErr.ErrorBadRequest
	}

	userId, err := controller.AuthService.GetValidUser(ctx)
	if err != nil {
		return err
	}

	err = controller.UserService.LinkEmail(ctx.UserContext(), userId, *emailReq)
	if err != nil {
		return err
	}

	return ctx.SendString("successfully linked email")
}

func (controller *UserController) UpdatePhone(ctx *fiber.Ctx) error {
	phoneReq := new(user_model.UpdatePhoneRequest)

	err := ctx.BodyParser(phoneReq)
	if err != nil {
		return customErr.ErrorBadRequest
	}

	userId, err := controller.AuthService.GetValidUser(ctx)
	if err != nil {
		return err
	}

	err = controller.UserService.LinkPhone(ctx.UserContext(), userId, *phoneReq)
	if err != nil {
		return err
	}

	return ctx.SendString("successfully linked phone")
}

func (controller *UserController) UpdateAccount(ctx *fiber.Ctx) error {
	updateAccReq := new(user_model.UpdateAccountRequest)

	err := ctx.BodyParser(updateAccReq)
	if err != nil {
		return customErr.ErrorBadRequest
	}

	userId, err := controller.AuthService.GetValidUser(ctx)
	if err != nil {
		return err
	}

	err = controller.UserService.UpdateAccount(ctx.UserContext(), userId, *updateAccReq)
	if err != nil {
		return err
	}

	return ctx.SendString("successfully update profile")
}
