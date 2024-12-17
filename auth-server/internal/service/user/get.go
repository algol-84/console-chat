package auth

import (
	"context"

	"github.com/algol-84/auth/internal/logger"
	"github.com/algol-84/auth/internal/model"
	"github.com/algol-84/auth/internal/repository"
	"go.uber.org/zap"
)

func (s *service) Get(ctx context.Context, id int64) (*model.User, error) {
	// Запрос юзера из кэша
	user, err := s.cacheRepository.Get(ctx, id)
	if err == model.ErrorUserNotFound {
		// Запрос юзера из базы
		user, err = s.authRepository.Get(ctx, &repository.Filter{ID: id})
		if err != nil {
			logger.Error("user not found", zap.String("error", err.Error()))
			return nil, model.ErrorUserNotFound
		}

		// Добавить юзера в кэш
		_, err = s.cacheRepository.Create(ctx, user)
		if err != nil {
			logger.Error("failed to create user in cache", zap.String("error", err.Error()))
			return nil, model.ErrorCacheInternal
		}

		return user, nil
	}

	return user, nil
}
