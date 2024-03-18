package user_model

import "time"

type Creator struct {
	UserId      int
	Name        string
	ImageUrl    string
	FriendCount int
	CreatedAt   time.Time
}
