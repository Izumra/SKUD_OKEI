package storage

import "errors"

var (
	ErrUserNotFound = errors.New("Пользователь с такими данными не зарегестрирован")
	ErrUserExist    = errors.New("Аккаунт с такими данными уже зарегестрирован")
)
