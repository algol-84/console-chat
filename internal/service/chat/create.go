package chat

import (
	"context"
	"log"

	"github.com/algol-84/chat-server/internal/model"
)

func (s *service) Create(ctx context.Context, chat *model.Chat) (int64, error) {
	log.Println(chat)
	id, err := s.chatRepository.Create(ctx, chat)
	if err != nil {
		return 0, nil
	}

	return id, nil
}
