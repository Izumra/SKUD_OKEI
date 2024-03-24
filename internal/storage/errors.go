package storage

import "errors"

var (
	ErrUserNotFound = errors.New("user of the system is not found")
	ErrUserExist    = errors.New("user with provided data already exists")
)
