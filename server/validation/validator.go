package validation

import (
	"gopkg.in/go-playground/validator.v9"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator(validator *validator.Validate) *CustomValidator {
	return &CustomValidator{validator: validator}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
