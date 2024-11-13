package service

import (
	"context"

	"github.com/algol-84/auth/internal/model"
)

// UserService интерфейс CRUD пользователей
type UserService interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error
}

// AuthService интерфейс сервиса авторизации
type AuthService interface {
	Login(ctx context.Context) error
}
