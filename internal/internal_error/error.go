package internal_error

import "errors"

var (
	ErrEmptyLoginOrPassword = errors.New("login or password is empty")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidPassword      = errors.New("invalid password")
)
