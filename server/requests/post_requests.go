package requests

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required" `
	Content string `json:"content" validate:"required"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" validate:"required" `
	Content string `json:"content" validate:"required"`
}