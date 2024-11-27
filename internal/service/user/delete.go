package auth

import (
	"context"
	"fmt"

	"github.com/algol-84/auth/internal/logger"
	"go.uber.org/zap"
)

func (s *service) Delete(ctx context.Context, id int64) error {
	// Удалить юзера из базы
	err1 := s.authRepository.Delete(ctx, id)
	// Удалить юзера из кэша
	err2 := s.cacheRepository.Delete(ctx, id)
	if err1 != nil || err2 != nil {
		logger.Error("failed to delete user", zap.Int64("id", id))
		return fmt.Errorf("repo error")
	}

	return nil
}
