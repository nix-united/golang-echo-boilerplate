package requests

import validation "github.com/go-ozzo/ozzo-validation/v4"

type BasicPost struct {
	Title   string `json:"title" validate:"required" example:"Echo"`
	Content string `json:"content" validate:"required" example:"Echo is nice!"`
}

func (bp BasicPost) Validate() error {
	return validation.ValidateStruct(&bp,
		validation.Field(&bp.Title, validation.Required),
		validation.Field(&bp.Content, validation.Required),
	)
}

type CreatePostRequest struct {
	BasicPost
}

type UpdatePostRequest struct {
	BasicPost
}
