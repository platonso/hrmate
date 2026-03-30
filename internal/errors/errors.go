package errors

import "errors"

var (
	// Form errors
	ErrFormNotFound        = errors.New("form not found")
	ErrFormInvalidStatus   = errors.New("invalid status")
	ErrFormAlreadyRejected = errors.New("form has already been rejected")
	ErrFormAlreadyApproved = errors.New("form has already been approved")

	// Username errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotActive      = errors.New("user account is not active")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrForbidden          = errors.New("forbidden")
)
