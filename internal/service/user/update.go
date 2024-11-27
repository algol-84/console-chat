package auth

import (
	"context"
	"errors"

	"github.com/algol-84/auth/internal/logger"
	"github.com/algol-84/auth/internal/model"
	"go.uber.org/zap"
)

func (s *service) Update(ctx context.Context, user *model.UserUpdate) error {
	// Проверить что запрос не пустой
	if !user.Name.Valid && !user.Email.Valid && !user.Role.Valid {
		return errors.New("update request is empty")
	}

	err := s.authRepository.Update(ctx, user)
	if err != nil {
		logger.Error("failed to update user", zap.String("error", err.Error()))
		return err
	}

	// Удалить юзера из кэша
	err = s.cacheRepository.Delete(ctx, user.ID)
	if err != nil {
		logger.Error("failed to delete user from cache", zap.String("error", err.Error()))
		return err
	}

	return nil
}
