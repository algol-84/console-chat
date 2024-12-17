package auth

import (
	"github.com/algol-84/auth/internal/config"
	"github.com/algol-84/auth/internal/repository"
	def "github.com/algol-84/auth/internal/service"
)

type service struct {
	tokenConfig      config.TokenConfig
	authRepository   repository.AuthRepository
	accessRepository repository.AccessRepository
}

// NewService конструктор сервисного слоя
func NewService(
	tokenConfig config.TokenConfig,
	authRepository repository.AuthRepository,
	accessRepository repository.AccessRepository,
) def.AccessService {
	return &service{
		tokenConfig:      tokenConfig,
		authRepository:   authRepository,
		accessRepository: accessRepository,
	}
}
