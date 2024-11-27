package auth

import (
	"context"
	"log"

	"github.com/algol-84/auth/internal/logger"
	"github.com/algol-84/auth/internal/model"
	"github.com/algol-84/auth/internal/utils"
	"go.uber.org/zap"
)

// GetAccessToken возвращает акцесс токен
func (s *service) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, []byte(s.tokenConfig.RefreshToken()))
	if err != nil {
		return "", model.ErrorRefreshToken
	}

	// Ищем пользователя в кэше для заполнения роли в клэйме
	user, err := s.cacheRepository.Get(ctx, claims.ID)
	if err != nil {
		logger.Error("failed to get user", zap.String("error", err.Error()))
		return "", model.ErrorUserNotFound
	}
	log.Println(user)

	accessToken, err := utils.GenerateToken(model.UserInfo{
		Username: claims.Username,
		Role:     user.Role,
	},
		[]byte(s.tokenConfig.AccessToken()),
		s.tokenConfig.AccessTokenExpiration(),
	)
	if err != nil {
		logger.Error("failed to generate access token", zap.String("error", err.Error()))
		return "", model.ErrorAccessToken
	}

	return accessToken, nil
}
