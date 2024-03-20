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
