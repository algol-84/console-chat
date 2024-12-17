package auth

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/algol-84/auth/internal/logger"
	"github.com/algol-84/auth/internal/model"
	"github.com/algol-84/auth/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

const (
	// authPrefix префикс добавляется к токену для идентификации используемого метода аутентификации.
	// В случае JWT принято добавлять Bearer
	authPrefix = "Bearer "
)

var accessibleRole map[string]string

// Check определяет уровень доступа пользователя
func (s *service) Check(ctx context.Context, endpoint string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("metadata is not provided")
		return model.ErrorAccessDenied
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		logger.Error("authorization header is not provided")
		return model.ErrorAccessDenied
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		logger.Error("invalid authorization header format")
		return model.ErrorAccessDenied
	}
	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	claims, err := utils.VerifyToken(accessToken, []byte(s.tokenConfig.AccessToken()))
	if err != nil {
		logger.Error("access token is invalid", zap.String("error", err.Error()))
		return model.ErrorAccessDenied
	}

	accessibleMap, err := s.accessibleRoles(ctx)
	if err != nil {
		logger.Error("failed to get accessible roles", zap.String("error", err.Error()))
		return model.ErrorAccessDenied
	}

	role, ok := accessibleMap[endpoint]
	if !ok {
		// Если роль не найдена по умолчанию политика доступа - запретить все
		logger.Error("endpoint not found, access denied")
		return model.ErrorAccessDenied
	}

	if role != claims.Role {
		logger.Error("insufficient user role, access denied")
		return model.ErrorAccessDenied
	}

	log.Printf("access granted for user %s to endpoint %s", claims.Username, endpoint)

	return nil
}

// Возвращает мапу с адресом эндпоинта и ролью, которая имеет доступ к нему
func (s *service) accessibleRoles(ctx context.Context) (map[string]string, error) {
	if accessibleRole == nil {

		accessibleRole = make(map[string]string)
		var err error
		// Лезем в базу за данными о доступных ролях для каждого эндпоинта
		// Можно кэшировать данные, чтобы не лезть в базу каждый раз
		accessibleRole, err = s.accessRepository.Get(ctx)
		if err != nil {
			return nil, errors.New("role not found")
		}

		return accessibleRole, nil
	}

	return accessibleRole, nil
}
