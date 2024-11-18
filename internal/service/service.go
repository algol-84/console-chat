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

// AuthService интерфейс сервиса аутентификации
type AuthService interface {
	Login(ctx context.Context, username string, password string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
}

// AccessService интерфейс сервиса авторизации
type AccessService interface {
	Check(ctx context.Context, endpoint string) error
}
