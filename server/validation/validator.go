package validation

import (
	v "gopkg.in/go-playground/validator.v9"
)

const minPasswordLength = 8

type CustomValidator struct {
	validator *v.Validate
}

func NewCustomValidator(validator *v.Validate) *CustomValidator {
	validator.RegisterValidation("password", ValidatePassword)

	return &CustomValidator{validator: validator}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func ValidatePassword(fl v.FieldLevel) bool {
	return len(fl.Field().String()) >= minPasswordLength
}
