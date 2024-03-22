package service

import (
	"context"
	"openidea-idea-social-media-app/customErr"
	friend_model "openidea-idea-social-media-app/models/friend"
	"openidea-idea-social-media-app/repository"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type FriendsService interface {
	AddFriends(ctx context.Context, userId int, request friend_model.FriendRequest) error
	RemoveFriends(ctx context.Context, userId int, request friend_model.FriendRequest) error
	GetAllFriends(ctx context.Context, userId int, filterRequest map[string]string) (friend_model.FriendsPagingResponse, error)
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

	uId, err := strconv.Atoi(request.UserId)
	if err != nil {
		return customErr.ErrorBadRequest
	}

	userFriend := friend_model.Friend{
		UserIdRequester: userId,
		UserIdAccepter:  uId,
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

	uId, err := strconv.Atoi(request.UserId)
	if err != nil {
		return customErr.ErrorBadRequest
	}

	userFriend := friend_model.Friend{
		UserIdRequester: userId,
		UserIdAccepter:  uId,
	}

	err = service.FriendsRepository.Delete(ctx, userFriend)
	if err != nil {
		return err
	}

	return nil
}

func (service *FriendsServiceImpl) GetAllFriends(ctx context.Context, userId int, filterRequest map[string]string) (friend_model.FriendsPagingResponse, error) {
	req, err := validateFilterQueryMap(filterRequest)
	if err != nil {
		return friend_model.FriendsPagingResponse{}, customErr.ErrorBadRequest
	}

	err = service.Validator.Struct(req)
	if err != nil {
		return friend_model.FriendsPagingResponse{}, customErr.ErrorBadRequest
	}

	res, err := service.FriendsRepository.GetAll(ctx, userId, req)
	if err != nil {
		return friend_model.FriendsPagingResponse{}, err
	}

	var responseData []friend_model.FriendDataResponse
	for _, v := range res.Data {
		data := friend_model.FriendDataResponse{
			UserId:      v.UserId,
			Name:        v.Name,
			ImageUrl:    v.ImageUrl.String,
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

func validateFilterQueryMap(req map[string]string) (friend_model.FilterFriends, error) {
	var filters friend_model.FilterFriends

	sortByVal, isSortByExists := req["sortBy"]
	if !isSortByExists && sortByVal == "" {
		filters.SortBy = "createdAt"
	} else {
		filters.SortBy = sortByVal
	}

	if isSortByExists && sortByVal == "" {
		return friend_model.FilterFriends{}, fiber.NewError(400, "false sortBy")
	}

	orderByVal, isOrderByExists := req["orderBy"]
	if !isOrderByExists && orderByVal == "" {
		filters.OrderBy = "desc"
	} else {
		filters.OrderBy = orderByVal
	}

	if isOrderByExists && orderByVal == "" {
		return friend_model.FilterFriends{}, fiber.NewError(400, "false orderBy")
	}

	userOnlyVal, isUserOnlyExists := req["onlyFriend"]
	if !isUserOnlyExists && userOnlyVal == "" {
		filters.UserOnly = false
	} else {
		resultUserOnly, err := strconv.ParseBool(userOnlyVal)
		if err != nil {
			return friend_model.FilterFriends{}, fiber.NewError(400, "false only users")
		}

		filters.UserOnly = resultUserOnly
	}

	if isUserOnlyExists && userOnlyVal == "" {
		return friend_model.FilterFriends{}, fiber.NewError(400, "false only user")
	}

	limitVal, isLimitExists := req["limit"]
	if !isLimitExists && limitVal == "" {
		filters.Limit = 5
	} else {
		resultLimit, err := strconv.Atoi(limitVal)
		if err != nil {
			return friend_model.FilterFriends{}, customErr.ErrorBadRequest
		}

		if resultLimit < 0 {
			return friend_model.FilterFriends{}, customErr.ErrorBadRequest
		}

		if resultLimit <= 5 {
			filters.Limit = 5
		} else {
			filters.Limit = resultLimit
		}
	}

	if isLimitExists && limitVal == "" {
		return friend_model.FilterFriends{}, customErr.ErrorBadRequest
	}

	offsetVal, isOffsetExists := req["offset"]
	if !isOffsetExists && offsetVal == "" {
		filters.Offset = 0
	} else {
		resultOffset, err := strconv.Atoi(offsetVal)
		if err != nil {
			return friend_model.FilterFriends{}, customErr.ErrorBadRequest
		}

		filters.Offset = resultOffset
	}

	if isOffsetExists && offsetVal == "" {
		return friend_model.FilterFriends{}, customErr.ErrorBadRequest
	}

	filters.Search = req["search"]

	return filters, nil
}
