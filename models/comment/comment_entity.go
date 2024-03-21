package comment_model

import (
	"database/sql"
)

type CommentCreator struct {
	UserId      int
	Name        string
	ImageUrl    sql.NullString
	FriendCount int
	CreatedAt   string
}

type Comment struct {
	CommentId int
	PostId    int
	UserId    int
	Comment   string
	Creator   CommentCreator
	CreatedAt string
}
