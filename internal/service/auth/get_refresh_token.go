package auth

import (
	"context"

	"github.com/algol-84/auth/internal/model"
	"github.com/algol-84/auth/internal/utils"
)

func (s *service) GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error) {
	// Проверяем токен и получаем стандартные клэймы
	// Проверка времени действия токена уже происходит под капотом
	claims, err := utils.VerifyToken(oldRefreshToken, []byte(s.tokenConfig.RefreshToken()))
	if err != nil {
		return "", model.ErrorRefreshToken
	}

	// Ищем пользователя в кэше для заполнения роли в клэйме
	user, err := s.cacheRepository.Get(ctx, claims.ID)
	if err != nil {
		return "", model.ErrorUserNotFound
	}

	refreshToken, err := utils.GenerateToken(model.UserInfo{
		Username: claims.Username,
		Role:     user.Role,
	},
		[]byte(s.tokenConfig.RefreshToken()),
		s.tokenConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
