package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrWrongPassword      = errors.New("wrong password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrValidationFailed   = errors.New("validation failed")
)
