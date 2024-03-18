package post_entity

import (
	comment_entity "openidea-idea-social-media-app/models/comment"
	user_model "openidea-idea-social-media-app/models/user"
	"time"
)

type Post struct {
	PostId    int
	PostHtml  string
	Tags      []string
	Comment   []comment_entity.Comment
	Creator   user_model.Creator
	CreatedAt time.Time
}
