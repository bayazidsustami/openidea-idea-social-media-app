package service

import (
	"context"
	"openidea-idea-social-media-app/customErr"
	comment_model "openidea-idea-social-media-app/models/comment"
	"openidea-idea-social-media-app/repository"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type CommentService interface {
	Create(ctx context.Context, userId string, request comment_model.CommentRequest) error
}

type CommentServiceImpl struct {
	Validator         *validator.Validate
	CommentRepository repository.CommentRepository
}

func NewCommentService(
	validator *validator.Validate,
	commentRepository repository.CommentRepository,
) CommentService {
	return &CommentServiceImpl{
		Validator:         validator,
		CommentRepository: commentRepository,
	}
}

func (service *CommentServiceImpl) Create(ctx context.Context, userId string, request comment_model.CommentRequest) error {
	err := service.Validator.Struct(request)
	if err != nil {
		return customErr.ErrorBadRequest
	}

	castedPostId, err := strconv.Atoi(request.PostId)
	if err != nil {
		return customErr.ErrorNotFound
	}

	comment := comment_model.Comment{
		PostId:  castedPostId,
		Comment: request.Comment,
		UserId:  userId,
	}

	err = service.CommentRepository.Create(ctx, comment)
	if err != nil {
		return err
	}

	return nil
}
