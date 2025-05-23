package models

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrPostNotFound = errors.New("post not found")
)
