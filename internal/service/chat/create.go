package chat

import (
	"context"

	"github.com/algol-84/chat-server/internal/model"
)

func (s *service) Create(ctx context.Context, chat *model.Chat) (int64, error) {
	id, errTx := s.chatRepository.Create(ctx, chat)
	if errTx != nil {
		return 0, errTx
	}

	return id, nil
}
