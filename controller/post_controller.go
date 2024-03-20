package controller

import (
	"openidea-idea-social-media-app/customErr"
	post_model "openidea-idea-social-media-app/models/post"
	"openidea-idea-social-media-app/service"

	"github.com/gofiber/fiber/v2"
)

type PostController struct {
	PostService service.PostService
	AuthService service.AuthService
}

func NewPostController(
	postService service.PostService,
	authService service.AuthService,
) PostController {
	return PostController{
		PostService: postService,
		AuthService: authService,
	}
}

func (controller *PostController) Create(ctx *fiber.Ctx) error {
	request := new(post_model.PostCreateRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		return customErr.ErrorBadRequest
	}

	userId, err := controller.AuthService.GetValidUser(ctx)
	if err != nil {
		return err
	}

	err = controller.PostService.Create(ctx.UserContext(), userId, *request)
	if err != nil {
		return err
	}

	return ctx.SendString("successfully add post")
}

func (controller *PostController) GetAll(ctx *fiber.Ctx) error {
	filterRequest := new(post_model.PostFilters)

	err := ctx.QueryParser(filterRequest)
	if err != nil {
		return customErr.ErrorBadRequest
	}

	userId, err := controller.AuthService.GetValidUser(ctx)
	if err != nil {
		return err
	}

	response, err := controller.PostService.GetAll(ctx.UserContext(), userId, *filterRequest)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}
