package models

import "errors"

const (
	ErrMsgKey = "error_description"
)

var (
	ErrTraceLayout            = "%s | error: %v"
	ErrUsernameIsAlreadyTaken = errors.New("username is already taken")
	ErrUserDoesNotExist       = errors.New("user not exist")
	ErrInvalidEmail           = errors.New("invalid email")
	ErrInvalidGender          = errors.New("invalid gender, available is: F/M/O")
	ErrInvalidAge             = errors.New("invalid age, the age must be greater than 1 and less than 150")
	ErrUserIsGone             = errors.New("user is gone")
	ErrDeleteDeletedUser      = errors.New("user has been deleted once")
)
