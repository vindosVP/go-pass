// Package storage consists the storage errors
package storage

import "errors"

var (
	// ErrUserAlreadyExists - error if user already exists
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrUserNotExist - error if user with provided email does not exist
	ErrUserNotExist = errors.New("user does not exist")
)
