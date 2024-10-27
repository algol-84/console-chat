package chat

import (
	"context"
)

func (s *service) Delete(ctx context.Context, id int64) error {
	errTx := s.chatRepository.Delete(ctx, id)
	if errTx != nil {
		return errTx
	}

	return errTx
}
