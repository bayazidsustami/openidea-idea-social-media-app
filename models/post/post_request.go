package post_model

type PostFilters struct {
	Limit     int    `json:"limit" validate:"number"`
	Offset    int    `json:"offset" validate:"number"`
	Search    string `json:"search" validate:"string"`
	SearchTag string `json:"searchTag" validate:"dive,string"`
}

type PostCreateRequest struct {
	PostHtml string   `json:"postInHtml" validate:"required,min=2,max=500"`
	Tags     []string `json:"tags" validate:"required,min=0"`
}