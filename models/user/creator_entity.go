package user_model

import (
	"database/sql"
	"time"
)

type Creator struct {
	UserId      int
	Name        string
	ImageUrl    sql.NullString
	FriendCount int
	CreatedAt   time.Time
}
