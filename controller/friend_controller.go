package controller

import (
	"openidea-idea-social-media-app/customErr"
	friend_model "openidea-idea-social-media-app/models/friend"
	"openidea-idea-social-media-app/service"

	"github.com/gofiber/fiber/v2"
)

type FriendController struct {
	FriendService service.FriendsService
	AuthService   service.AuthService
}

func NewFriendsController(
	friendService service.FriendsService,
	authService service.AuthService,
) FriendController {
	return FriendController{
		FriendService: friendService,
		AuthService:   authService,
	}
}

func (controller *FriendController) AddFriend(ctx *fiber.Ctx) error {
	request := new(friend_model.FriendRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		return customErr.ErrorBadRequest
	}

	userId, err := controller.AuthService.GetValidUser(ctx)
	if err != nil {
		return err
	}

	err = controller.FriendService.AddFriends(ctx.UserContext(), userId, *request)
	if err != nil {
		return err
	}

	return ctx.SendString("successfully add friend")
}

func (controller *FriendController) RemoveFriends(ctx *fiber.Ctx) error {
	request := new(friend_model.FriendRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		return customErr.ErrorBadRequest
	}

	userId, err := controller.AuthService.GetValidUser(ctx)
	if err != nil {
		return err
	}

	err = controller.FriendService.RemoveFriends(ctx.UserContext(), userId, *request)
	if err != nil {
		return err
	}

	return ctx.SendString("successfully delete friend")
}
