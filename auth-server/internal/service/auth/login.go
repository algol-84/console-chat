package auth

import (
	"context"

	"github.com/algol-84/auth/internal/logger"
	"github.com/algol-84/auth/internal/model"
	"github.com/algol-84/auth/internal/repository"
	"github.com/algol-84/auth/internal/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Login(ctx context.Context, username string, password string) (string, error) {
	// Найти юзера по имени в базе
	user, err := s.authRepository.Get(ctx, &repository.Filter{Username: username})
	if err != nil {
		logger.Error("failed to get user", zap.String("error", err.Error()))
		return "", model.ErrorRefreshToken
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.Error("failed to compare password", zap.String("error", err.Error()))
		return "", model.ErrorRefreshToken
	}

	// Добавить юзера в кэш
	_, err = s.cacheRepository.Create(ctx, user)
	if err != nil {
		logger.Error("failed to create user in cache", zap.String("error", err.Error()))
		return "", model.ErrorCacheInternal
	}

	// Генерируем токен для конкретного юзера из базы
	refreshToken, err := utils.GenerateToken(model.UserInfo{
		ID:       user.ID,
		Username: user.Name,
		Role:     user.Role,
	},
		[]byte(s.tokenConfig.RefreshToken()),
		s.tokenConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		logger.Error("failed to generate token", zap.String("error", err.Error()))
		return "", model.ErrorRefreshToken
	}

	return refreshToken, nil
}
