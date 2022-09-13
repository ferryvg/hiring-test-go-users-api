package dal

import "errors"

var (
	ErrNotFound           = errors.New("entity not found")
	ErrUserExists         = errors.New("user with provided id already exists")
	ErrInvalidCredentials = errors.New("invalid user credentials")
)
