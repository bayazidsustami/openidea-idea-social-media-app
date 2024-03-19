package post_model

type PostCreateRequest struct {
	PostHtml string   `json:"postInHtml" validate:"required,min=2,max=500"`
	Tags     []string `json:"tags" validate:"required,min=0"`
}
