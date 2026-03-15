package errors

import "errors"

var (
	// Form errors
	ErrFormNotFound = errors.New("form not found")

	// Username errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrForbidden          = errors.New("forbidden")
)
