package app

import (
	"context"
	"log"

	auth "github.com/algol-84/auth/internal/api/auth"
	"github.com/algol-84/auth/internal/client/db"
	"github.com/algol-84/auth/internal/client/db/pg"
	"github.com/algol-84/auth/internal/closer"
	"github.com/algol-84/auth/internal/config"
	"github.com/algol-84/auth/internal/repository"
	authRepository "github.com/algol-84/auth/internal/repository/auth"
	"github.com/algol-84/auth/internal/service"
	authService "github.com/algol-84/auth/internal/service/auth"
	// "github.com/jackc/pgx/v4/pgxpool"
)

// serviceProvider хранит все объекты приложения, как интерфейсы или ссылки на структуры
type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	authRepository repository.AuthRepository

	authService service.AuthService

	authImpl *auth.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// Определяются функции инициализации всех объектов

// PGConfig инициализирует считывание настроек PG из файла конфига
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

// GRPCConfig инициализирует считывание настроек GRPC из файла конфига
func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepository.NewRepository(s.DBClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(s.AuthRepository(ctx))
	}

	return s.authService
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}
