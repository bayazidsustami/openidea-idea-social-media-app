package comment_model

type CommentRequest struct {
	PostId  string `json:"postId" validate:"required"`
	Comment string `json:"comment" validate:"required,min=2,max=500"`
}
