package comment_model

type CommentRequest struct {
	// TODO: Add validation 'should be a valid post id'
	PostId  string `json:"postId" validate:"required"`
	Comment string `json:"comment" validate:"required,min=2,max=500"`
}
