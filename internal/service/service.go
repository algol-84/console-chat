package service

import (
	"context"

	"github.com/algol-84/chat-server/internal/model"
)

// ChatService интерфейс
type ChatService interface {
	Create(ctx context.Context, user *model.Chat) (int64, error)
	Delete(ctx context.Context, id int64) error
}
