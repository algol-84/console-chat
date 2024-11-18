package auth

import (
	"github.com/algol-84/auth/internal/client/kafka"
	"github.com/algol-84/auth/internal/repository"
	def "github.com/algol-84/auth/internal/service"
)

type service struct {
	authRepository  repository.AuthRepository
	cacheRepository repository.CacheRepository
	kafkaProducer   kafka.Producer
}

// NewService конструктор сервисного слоя
func NewService(
	authRepository repository.AuthRepository,
	cacheRepository repository.CacheRepository,
	kafkaProducer kafka.Producer,
) def.UserService {
	return &service{
		authRepository:  authRepository,
		cacheRepository: cacheRepository,
		kafkaProducer:   kafkaProducer,
	}
}

// NewMockService Mock-конструктор сервисного слоя
func NewMockService(deps ...interface{}) def.UserService {
	srv := service{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.AuthRepository:
			srv.authRepository = s
		case repository.CacheRepository:
			srv.cacheRepository = s
		case kafka.Producer:
			srv.kafkaProducer = s
		}
	}

	return &srv
}
