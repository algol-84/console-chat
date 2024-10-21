package auth

import (
	"github.com/algol-84/auth/internal/service"
	desc "github.com/algol-84/auth/pkg/user_v1"
)

// Implementation содержит все объекты сервера
type Implementation struct {
	desc.UnimplementedUserV1Server
	authService service.AuthService
}

// NewImplementation конструктор API слоя
func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
