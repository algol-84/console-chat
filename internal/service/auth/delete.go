package auth

import (
	"context"
)

func (s *service) Delete(ctx context.Context, id int64) error {
	err := s.authRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
