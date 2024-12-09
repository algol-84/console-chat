package model

import (
	"time"
)

// User представляет модель пользователя сервисного слоя
type User struct {
	ID              int64
	Name            string
	Password        string
	PasswordConfirm string
	Email           string
	Role            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
