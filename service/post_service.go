package service

import (
	"context"
	"openidea-idea-social-media-app/customErr"
	common_model "openidea-idea-social-media-app/models/common"
	post_model "openidea-idea-social-media-app/models/post"
	"openidea-idea-social-media-app/repository"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type PostService interface {
	Create(ctx context.Context, userId string, request post_model.PostCreateRequest) error
	GetAll(ctx context.Context, filters map[string]string) (post_model.PostGetAllResponse, error)
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

func (service *PostServiceImpl) Create(ctx context.Context, userId string, request post_model.PostCreateRequest) error {
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

func (service *PostServiceImpl) GetAll(ctx context.Context, request map[string]string) (post_model.PostGetAllResponse, error) {
	filters, err := validateQueryParams(request)
	if err != nil {
		return post_model.PostGetAllResponse{}, err
	}
	err = service.Validator.Struct(filters)
	if err != nil {
		return post_model.PostGetAllResponse{}, customErr.ErrorBadRequest
	}

	posts, totalPosts, err := service.PostRepository.GetAll(ctx, filters)
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
					ImageUrl:    comment.Creator.ImageUrl,
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
		Meta: common_model.MetaPageResponse{
			Limit:  filters.Limit,
			Offset: filters.Offset,
			Total:  totalPosts,
		},
	}

	return response, nil
}

func validateQueryParams(req map[string]string) (post_model.PostFilters, error) {
	var filters post_model.PostFilters
	limitVal, isLimitExists := req["limit"]
	if !isLimitExists && limitVal == "" {
		filters.Limit = 5
	} else {
		resultLimit, err := strconv.Atoi(limitVal)
		if err != nil {
			return post_model.PostFilters{}, customErr.ErrorBadRequest
		}

		if resultLimit < 0 {
			return post_model.PostFilters{}, customErr.ErrorBadRequest
		}

		if resultLimit <= 5 {
			filters.Limit = 5
		} else {
			filters.Limit = resultLimit
		}
	}

	if isLimitExists && limitVal == "" {
		return post_model.PostFilters{}, customErr.ErrorBadRequest
	}

	offsetVal, isOffsetExists := req["offset"]
	if !isOffsetExists && offsetVal == "" {
		filters.Offset = 0
	} else {
		resultOffset, err := strconv.Atoi(offsetVal)
		if err != nil {
			return post_model.PostFilters{}, customErr.ErrorBadRequest
		}

		filters.Offset = resultOffset
	}

	if isOffsetExists && offsetVal == "" {
		return post_model.PostFilters{}, customErr.ErrorBadRequest
	}

	filters.Search = req["search"]
	filters.SearchTag = strings.Split(req["searchTag"], ",")

	return filters, nil
}
