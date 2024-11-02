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
	Name  StringValue
	Email StringValue
	Role  StringValue
}

// StringValue кастомный тип строки, если Valid=false, то строка не валидна
type StringValue struct {
	Value string
	Valid bool
}
