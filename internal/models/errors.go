package models

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrInvalidAuthToken = errors.New("invalid authorization jwt token")

	ErrPostNotFound = errors.New("post not found")
)
