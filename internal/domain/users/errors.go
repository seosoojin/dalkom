package users

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound           = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidUsername    = errors.New("invalid username")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidPassword    = errors.New("invalid password")
)
