package chat

import (
	"context"

	"github.com/algol-84/chat-server/internal/model"
)

func (s *service) Create(ctx context.Context, chat *model.Chat) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.chatRepository.Create(ctx, chat)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.logRepository.Create(ctx, id, "chat was created")
		if errTx != nil {
			return errTx
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
