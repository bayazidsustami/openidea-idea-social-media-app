package service

import (
	"context"
	"openidea-idea-social-media-app/customErr"
	common_model "openidea-idea-social-media-app/models/common"
	post_model "openidea-idea-social-media-app/models/post"
	"openidea-idea-social-media-app/repository"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type PostService interface {
	Create(ctx context.Context, userId int, request post_model.PostCreateRequest) error
	GetAll(ctx context.Context, userId int, filters post_model.PostFilters) (post_model.PostGetAllResponse, error)
}

type PostServiceImpl struct {
	Validator      *validator.Validate
	PostRepository repository.PostRepository
}

func NewPostService(
	validator *validator.Validate,
	postRepository repository.PostRepository,
) PostService {
	return &PostServiceImpl{
		Validator:      validator,
		PostRepository: postRepository,
	}
}

func (service *PostServiceImpl) Create(ctx context.Context, userId int, request post_model.PostCreateRequest) error {
	err := service.Validator.Struct(request)
	if err != nil {
		return customErr.ErrorBadRequest
	}

	post := post_model.Post{
		PostHtml: request.PostHtml,
		Tags:     request.Tags,
	}

	err = service.PostRepository.Create(ctx, post, userId)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Populate userId
func (service *PostServiceImpl) GetAll(ctx context.Context, userId int, filters post_model.PostFilters) (post_model.PostGetAllResponse, error) {
	err := service.Validator.Struct(filters)
	if err != nil {
		return post_model.PostGetAllResponse{}, customErr.ErrorBadRequest
	}

	posts, err := service.PostRepository.GetAll(ctx, filters)
	if err != nil {
		return post_model.PostGetAllResponse{}, err
	}

	var data []post_model.PostDataResponse
	for _, post := range posts {
		var comments []post_model.CommentResponse
		for _, comment := range post.Comments {
			comments = append(comments, post_model.CommentResponse{
				Comment:   comment.Comment,
				CreatedAt: comment.CreatedAt,
				Creator: post_model.CreatorResponse{
					UserId:      comment.Creator.UserId,
					Name:        comment.Creator.Name,
					ImageUrl:    comment.Creator.ImageUrl.String,
					FriendCount: comment.Creator.FriendCount,
					CreatedAt:   comment.Creator.CreatedAt,
				},
			})

		}

		rowData := post_model.PostDataResponse{
			PostId: strconv.Itoa(post.PostId),
			Post: post_model.PostResponse{
				PostInHtml: post.PostHtml,
				Tags:       post.Tags,
				CreatedAt:  post.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			},
			Comments: comments,
			Creator: post_model.CreatorResponse{
				UserId:      post.Creator.UserId,
				Name:        post.Creator.Name,
				ImageUrl:    post.Creator.ImageUrl.String,
				FriendCount: post.Creator.FriendCount,
				CreatedAt:   post.Creator.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			},
		}

		data = append(data, rowData)
	}

	response := post_model.PostGetAllResponse{
		Message: "Success",
		Data:    data,
		// TODO: Populate this with actual value
		Meta: common_model.MetaPageResponse{},
	}

	return response, nil
}
