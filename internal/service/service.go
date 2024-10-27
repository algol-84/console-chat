package service

import (
	"context"

	"github.com/algol-84/auth/internal/model"
)

// AuthService интерфейс
type AuthService interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error
}
