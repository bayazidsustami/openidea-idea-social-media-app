package friend_model

import (
	"database/sql"
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
	ImageUrl    sql.NullString
	FriendCount int
	CreatedAt   *time.Time
}

type FriendDataPaging struct {
	Data []FriendData
	Meta models.MetaPage
}
