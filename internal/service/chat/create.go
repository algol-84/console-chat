package chat

import (
	"context"
)

func (s *service) Create(_ context.Context) (int64, error) {
	// id, errTx := s.chatRepository.Create(ctx, chat)
	// if errTx != nil {
	// 	return 0, errTx
	// }

	id := int64(1)

	return id, nil
}
