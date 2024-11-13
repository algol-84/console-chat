package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/algol-84/auth/internal/model"
	"github.com/algol-84/auth/internal/utils"
)

// TODO вынести константы в конфиг
const (
	authPrefix = "Bearer "

	// TODO генерировать токены через TLS
	refreshTokenSecretKey = "W4/X+LLjehdxptt4YgGFCvMpq5ewptpZZYRHY6A72g0="
	accessTokenSecretKey  = "VqvguGiffXILza1f44TWXowDT4zwf03dtXmqWW4SYyE="

	refreshTokenExpiration = 60 * time.Minute
	accessTokenExpiration  = 5 * time.Minute
)

func (s *service) Login(ctx context.Context, username string, password string) (string, error) {
	log.Println("login handle in service layer", username, password)
	// Найти юзера по имени в базе
	user, err := s.authRepository.Find(ctx, username)
	if err != nil {
		return "", errors.New("user not found")
	}

	log.Println(user)
	// Генерируем токен для конкретного юзера из базы
	refreshToken, err := utils.GenerateToken(model.UserInfo{
		Username: user.Name,
		Role: user.Role,
	},
		[]byte(refreshTokenSecretKey),
		refreshTokenExpiration,
	)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return refreshToken, nil
}
