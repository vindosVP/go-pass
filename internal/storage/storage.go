package storage

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotExist      = errors.New("user does not exist")
)
