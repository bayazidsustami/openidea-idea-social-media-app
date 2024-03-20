package service

import (
	"context"
	"openidea-idea-social-media-app/customErr"
	friend_model "openidea-idea-social-media-app/models/friend"
	"openidea-idea-social-media-app/repository"

	"github.com/go-playground/validator/v10"
)

type FriendsService interface {
	AddFriends(ctx context.Context, userId int, request friend_model.FriendRequest) error
	RemoveFriends(ctx context.Context, userId int, request friend_model.FriendRequest) error
	GetAllFriends(ctx context.Context, userId int, filterRequest friend_model.FilterFriends) (friend_model.FriendsPagingResponse, error)
}

type FriendsServiceImpl struct {
	Validator         *validator.Validate
	FriendsRepository repository.FriendsRepository
}

func NewFriendsService(
	validator *validator.Validate,
	friendRepository repository.FriendsRepository,
) FriendsService {
	return &FriendsServiceImpl{
		Validator:         validator,
		FriendsRepository: friendRepository,
	}
}

func (service *FriendsServiceImpl) AddFriends(ctx context.Context, userId int, request friend_model.FriendRequest) error {
	err := service.Validator.Struct(request)
	if err != nil {
		return customErr.ErrorBadRequest
	}

	userFriend := friend_model.Friend{
		UserIdRequester: userId,
		UserIdAccepter:  request.UserId,
	}

	err = service.FriendsRepository.Create(ctx, userFriend)
	if err != nil {
		return err
	}

	return nil
}

func (service *FriendsServiceImpl) RemoveFriends(ctx context.Context, userId int, request friend_model.FriendRequest) error {
	err := service.Validator.Struct(request)
	if err != nil {
		return customErr.ErrorBadRequest
	}

	userFriend := friend_model.Friend{
		UserIdRequester: userId,
		UserIdAccepter:  request.UserId,
	}

	err = service.FriendsRepository.Delete(ctx, userFriend)
	if err != nil {
		return err
	}

	return nil
}

func (service *FriendsServiceImpl) GetAllFriends(ctx context.Context, userId int, filterRequest friend_model.FilterFriends) (friend_model.FriendsPagingResponse, error) {
	err := service.Validator.Struct(filterRequest)
	if err != nil {
		return friend_model.FriendsPagingResponse{}, customErr.ErrorBadRequest
	}

	if filterRequest.SortBy == "" {
		filterRequest.SortBy = "createdAt"
	}

	if filterRequest.OrderBy == "" {
		filterRequest.OrderBy = "desc"
	}

	if filterRequest.Limit <= 5 {
		filterRequest.Limit = 5
	}

	if filterRequest.Offset == 0 {
		filterRequest.Offset = 0
	}

	res, err := service.FriendsRepository.GetAll(ctx, userId, filterRequest)
	if err != nil {
		return friend_model.FriendsPagingResponse{}, err
	}

	var responseData []friend_model.FriendDataResponse
	for _, v := range res.Data {
		data := friend_model.FriendDataResponse{
			UserId:      v.UserId,
			Name:        v.Name,
			ImageUrl:    v.ImageUrl,
			FriendCount: v.FriendCount,
			CreatedAt:   v.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
		}
		responseData = append(responseData, data)
	}

	return friend_model.FriendsPagingResponse{
		Message: "Success",
		Data:    responseData,
		Meta:    res.Meta,
	}, nil
}
