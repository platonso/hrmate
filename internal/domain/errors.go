package domain

import "errors"

var (
	ErrFormNotFound = errors.New("form not found")

	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
