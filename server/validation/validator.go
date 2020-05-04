package validation

import (
	v "gopkg.in/go-playground/validator.v9"
)

type CustomValidator struct {
	validator *v.Validate
}

func NewCustomValidator(validator *v.Validate) *CustomValidator {
	return &CustomValidator{validator: validator}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
