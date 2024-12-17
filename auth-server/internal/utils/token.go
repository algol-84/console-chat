package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/algol-84/auth/internal/model"
)

// GenerateToken генерирует JWT токен, подписанный секретным ключом secretKey и валидный duration секунд
func GenerateToken(info model.UserInfo, secretKey []byte, duration time.Duration) (string, error) {
	// claims содержит payload jwt токена
	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			// Инициализация стандартных клэймов
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		ID:       info.ID,
		Username: info.Username,
		Role:     info.Role,
	}

	// Создаем токен на основе выбранного метода криптографии и передаем клэймы
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Подписываем токен секретным ключом и возвращаем
	return token.SignedString(secretKey)
}

// VerifyToken проверяет токен с помощью секретного ключа и возвращает его клэймы
func VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.Errorf("unexpected token signing method")
			}

			return secretKey, nil
		},
	)
	if err != nil {
		return nil, errors.Errorf("invalid token: %s", err.Error())
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, errors.Errorf("invalid token claims")
	}

	return claims, nil
}
