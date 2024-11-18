package auth

import (
	"context"
	"fmt"
)

func (s *service) Delete(ctx context.Context, id int64) error {
	// Удалить юзера из базы
	err1 := s.authRepository.Delete(ctx, id)
	// Удалить юзера из кэша
	err2 := s.cacheRepository.Delete(ctx, id)
	if err1 != nil || err2 != nil {
		return fmt.Errorf("repo error")
	}

	return nil
}
