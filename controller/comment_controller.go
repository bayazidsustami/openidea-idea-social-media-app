package controller

import (
	"openidea-idea-social-media-app/customErr"
	comment_model "openidea-idea-social-media-app/models/comment"
	"openidea-idea-social-media-app/service"

	"github.com/gofiber/fiber/v2"
)

type CommentController struct {
	CommentService service.CommentService
	AuthService    service.AuthService
}

func NewCommentController(
	commentService service.CommentService,
	authService service.AuthService,
) CommentController {
	return CommentController{
		CommentService: commentService,
		AuthService:    authService,
	}
}

func (controller *CommentController) Create(ctx *fiber.Ctx) error {
	request := new(comment_model.CommentRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		return customErr.ErrorBadRequest
	}

	userId, err := controller.AuthService.GetValidUser(ctx)
	if err != nil {
		return err
	}

	err = controller.CommentService.Create(ctx.UserContext(), userId, *request)
	if err != nil {
		return err
	}

	return ctx.SendString("successfully add comment")
}
