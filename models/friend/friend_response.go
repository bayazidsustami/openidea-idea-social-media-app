package friend_model

import "openidea-idea-social-media-app/models"

type FriendsPagingResponse struct {
	Message string               `json:"message"`
	Data    []FriendDataResponse `json:"data"`
	Meta    models.MetaPage      `json:"meta"`
}

type FriendDataResponse struct {
	UserId      int    `json:"userId"`
	Name        string `json:"name"`
	ImageUrl    string `json:"imageUrl"`
	FriendCount int    `json:"friendCount"`
	CreatedAt   string `json:"createdAt"`
}
