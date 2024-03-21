package post_model

import (
	common_model "openidea-idea-social-media-app/models/common"
)

type CreatorResponse struct {
	UserId      int    `json:"userId"`
	Name        string `json:"name"`
	ImageUrl    string `json:"imageUrl"`
	FriendCount int    `json:"friendCount"`
	CreatedAt   string `json:"createdAt"`
}

type PostResponse struct {
	PostInHtml string   `json:"postInHtml"`
	Tags       []string `json:"tags"`
	CreatedAt  string   `json:"createdAt"`
}

type CommentResponse struct {
	Comment   string          `json:"comment"`
	Creator   CreatorResponse `json:"creator"`
	CreatedAt string          `json:"createdAt"`
}

type PostDataResponse struct {
	PostId   string            `json:"postId"`
	Post     PostResponse      `json:"post"`
	Comments []CommentResponse `json:"comments"`
	Creator  CreatorResponse   `json:"creator"`
}

type PostGetAllResponse struct {
	Message string                        `json:"message"`
	Data    []PostDataResponse            `json:"data"`
	Meta    common_model.MetaPageResponse `json:"meta"`
}
