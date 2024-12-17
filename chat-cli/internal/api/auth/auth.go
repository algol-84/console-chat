package auth

import (
	"context"

	"github.com/algol-84/chat-cli/pkg/auth_v1"
)

// Login осуществляет логин пользователя в сервисе авторизации
// Если такого пользователя не существует в базе, то возвращается ошибка
// в случае успеха возвращается JWT Refresh Token
func (a *AuthImpl) Login(username string, password string) (string, error) {
	ctx := context.Background()
	res, err := a.authClient.Login(ctx, &auth_v1.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", err
	}

	return res.RefreshToken, nil
}

// GetAccessToken делает запрос на выпуск JWT access токена
func (a *AuthImpl) GetAccessToken(refreshToken string) (string, error) {
	ctx := context.Background()
	res, err := a.authClient.GetAccessToken(ctx, &auth_v1.GetAccessTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return "", err
	}

	return res.AccessToken, nil
}
