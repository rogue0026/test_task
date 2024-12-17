package storage

import "errors"

var (
	ErrLoginAlreadyExists = errors.New("login already exists")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
)
