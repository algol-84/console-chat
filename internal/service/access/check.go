package auth

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/algol-84/auth/internal/utils"
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
	log.Println("111")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("metadata is not provided")
	}
	log.Println(md)
	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return errors.New("invalid authorization header format")
	}
	log.Println(authHeader[0])
	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	claims, err := utils.VerifyToken(accessToken, []byte(s.tokenConfig.AccessToken()))
	if err != nil {
		log.Println("access token is invalid")
		return errors.New("access token is invalid")
	}

	log.Println("333")
	accessibleMap, err := s.accessibleRoles(ctx)
	if err != nil {
		return errors.New("failed to get accessible roles")
	}

	log.Println(accessibleMap)

	role, ok := accessibleMap[endpoint]
	if !ok {
		// Если роль не найдена по умолчанию политика доступа - запретить все
		return errors.New("endpoint not found, access denied")
	}

	log.Println(role, claims.Role)
	if role != claims.Role {
		return errors.New("access denied")
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
