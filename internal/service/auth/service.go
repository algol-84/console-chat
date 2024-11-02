package auth

import (
	"github.com/algol-84/auth/internal/repository"
	def "github.com/algol-84/auth/internal/service"
)

type service struct {
	authRepository  repository.AuthRepository
	cacheRepository repository.CacheRepository
}

// NewService конструктор сервисного слоя
func NewService(authRepository repository.AuthRepository, cacheRepository repository.CacheRepository) def.AuthService {
	return &service{
		authRepository:  authRepository,
		cacheRepository: cacheRepository,
	}
}

// NewMockService Mock-конструктор сервисного слоя
func NewMockService(deps ...interface{}) def.AuthService {
	srv := service{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.AuthRepository:
			srv.authRepository = s
		case repository.CacheRepository:
			srv.cacheRepository = s
		}
	}

	return &srv
}
