package auth

import (
	"github.com/algol-84/auth/internal/service"
	descAccess "github.com/algol-84/auth/pkg/access_v1"
)

// Implementation содержит все объекты сервера
type Implementation struct {
	descAccess.UnimplementedAccessV1Server
	accessService service.AccessService
}

// NewImplementation конструктор API слоя
func NewImplementation(accessService service.AccessService) *Implementation {
	return &Implementation{
		accessService: accessService,
	}
}
