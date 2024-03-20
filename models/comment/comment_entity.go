package comment_model

import user_model "openidea-idea-social-media-app/models/user"

type Comment struct {
	CommentId int
	PostId    int
	UserId    int
	Comment   string
	Creator   user_model.Creator
}
