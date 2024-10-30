package model

import "github.com/pkg/errors"

// ErrorUserNotFound ошибка пользователь не найден
var ErrorUserNotFound = errors.New("user not found")
