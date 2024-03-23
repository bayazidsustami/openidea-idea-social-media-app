package friend_model

import (
	"database/sql"
	"openidea-idea-social-media-app/models"
	"time"
)

type Friend struct {
	UserIdRequester string
	UserIdAccepter  string
}

type FriendData struct {
	UserId      string
	Name        string
	ImageUrl    sql.NullString
	FriendCount int
	CreatedAt   *time.Time
}

type FriendDataPaging struct {
	Data []FriendData
	Meta models.MetaPage
}
