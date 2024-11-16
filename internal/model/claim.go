package model

import "github.com/dgrijalva/jwt-go"

// UserClaims структура кастомных клэймов
type UserClaims struct {
	// Стандартная структура JWT Claims из RFC
	jwt.StandardClaims
	// Наши кастомные поля, которые мы хотим встроить
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// UserInfo модель юзера, которая содержит все данные для сохранения в payload jwt токена
type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
