package post_model

import (
	comment_model "openidea-idea-social-media-app/models/comment"
	user_model "openidea-idea-social-media-app/models/user"
	"time"
)

type Post struct {
	PostId    int
	PostHtml  string
	Tags      []string
	Comments  []comment_model.Comment
	Creator   user_model.Creator
	CreatedAt time.Time
}
