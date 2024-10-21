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

// UserUpdate представляет модель апдейта юзера сервисного слоя
type UserUpdate struct {
	ID    int64
	Name  stringValue
	Email stringValue
	Role  stringValue
}

type stringValue struct {
	Value string
	Valid bool
}
