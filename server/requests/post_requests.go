package requests

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required" example:"Echo"`
	Content string `json:"content" validate:"required" example:"Echo is nice!"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" validate:"required" example:"Echo"`
	Content string `json:"content" validate:"required" example:"Echo is very nice!"`
}
