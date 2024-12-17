package model

import (
	"database/sql"
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
	UpdatedAt       sql.NullTime
}
