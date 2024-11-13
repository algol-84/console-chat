package auth

import (
	"github.com/algol-84/auth/internal/service"
	descUser "github.com/algol-84/auth/pkg/user_v1"
)

// Implementation содержит все объекты сервера
type Implementation struct {
	descUser.UnimplementedUserV1Server
	userService service.UserService
}

// NewImplementation конструктор API слоя
func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
