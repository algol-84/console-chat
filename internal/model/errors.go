package model

import "github.com/pkg/errors"

// ErrorUserNotFound ошибка пользователь не найден
var ErrorUserNotFound = errors.New("user not found")

// ErrorCacheInternal внутрення ошибка в базе Redis
var ErrorCacheInternal = errors.New("internal cache error")

// ErrorRefreshToken неверный рефреш токен
var ErrorRefreshToken = errors.New("invalid refresh token")