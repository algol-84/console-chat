package auth

import (
	"github.com/algol-84/auth/internal/service"
	descAuth "github.com/algol-84/auth/pkg/auth_v1"
)

// Implementation содержит все объекты сервера
type Implementation struct {
	descAuth.UnimplementedAuthV1Server
	authService service.AuthService
}

// NewImplementation конструктор API слоя
func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
