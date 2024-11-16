package auth

import (
	"context"
	"errors"

	"github.com/algol-84/auth/internal/model"
	"github.com/algol-84/auth/internal/utils"
)

func (s *service) Login(ctx context.Context, username string, password string) (string, error) {
	// Найти юзера по имени в базе
	user, err := s.authRepository.Find(ctx, username, password)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Добавить юзера в кэш
	_, err = s.cacheRepository.Create(ctx, user)
	if err != nil {
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
		return "", errors.New("failed to generate token")
	}

	return refreshToken, nil
}
