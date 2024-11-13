package model

import "github.com/dgrijalva/jwt-go"

const (
	ExamplePath = "/note_v1.NoteV1/Get"
)

type UserClaims struct {
	// Стандартная структура JWT Claims из RFC
	jwt.StandardClaims		
	// Наши кастомные поля, которые мы хотим встроить
	Username string `json:"username"`
	Role     string `json:"role"`
}
