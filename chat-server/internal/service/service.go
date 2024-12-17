package service

import (
	"context"
)

// ChatService интерфейс
type ChatService interface {
	Create(ctx context.Context) (int64, error)
	Delete(ctx context.Context, id int64) error
}
