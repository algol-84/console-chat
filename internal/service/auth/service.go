package auth

import (
	"github.com/algol-84/auth/internal/repository"
	def "github.com/algol-84/auth/internal/service"
)

type service struct {
	authRepository repository.AuthRepository
}

// NewService конструктор сервисного слоя
func NewService(
	authRepository repository.AuthRepository,
) def.AuthService {
	return &service{
		authRepository: authRepository,
	}
}
