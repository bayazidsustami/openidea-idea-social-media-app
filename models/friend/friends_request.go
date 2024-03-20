package friend_model

type FriendRequest struct {
	UserId int `json:"userId" validate:"required"`
}
