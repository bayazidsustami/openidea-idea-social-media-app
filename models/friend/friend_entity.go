package friend_model

import (
	"openidea-idea-social-media-app/models"
	"time"
)

type Friend struct {
	UserIdRequester int
	UserIdAccepter  int
}

type FriendData struct {
	UserId      int
	Name        string
	ImageUrl    string
	FriendCount string
	CreatedAt   *time.Time
}

type FriendDataPaging struct {
	Data []FriendData
	Meta models.MetaPage
}
