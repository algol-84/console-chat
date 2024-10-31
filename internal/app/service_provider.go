package app

import (
	"context"
	"log"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/algol-84/auth/internal/api/auth"
	"github.com/algol-84/auth/internal/client/cache"
	"github.com/algol-84/auth/internal/client/cache/redis"
	"github.com/algol-84/auth/internal/config"
	"github.com/algol-84/auth/internal/repository"
	authRepositoryPg "github.com/algol-84/auth/internal/repository/auth/pg"
	authRepositoryRedis "github.com/algol-84/auth/internal/repository/auth/redis"
	"github.com/algol-84/auth/internal/service"
	authService "github.com/algol-84/auth/internal/service/auth"
	closer "github.com/algol-84/platform_common/pkg/closer"
	db "github.com/algol-84/platform_common/pkg/db"
	pg "github.com/algol-84/platform_common/pkg/db/pg"
)

// serviceProvider хранит все объекты приложения, как интерфейсы или ссылки на структуры
type serviceProvider struct {
	pgConfig    config.PGConfig
	grpcConfig  config.GRPCConfig
	redisConfig config.RedisConfig

	dbClient        db.Client
	redisPool       *redigo.Pool
	redisClient     cache.RedisClient
	authRepository  repository.AuthRepository
	cacheRepository repository.CacheRepository

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

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
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

func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.RedisConfig().MaxIdle(),
			IdleTimeout: s.RedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", s.RedisConfig().Address())
			},
		}
	}

	return s.redisPool
}

func (s *serviceProvider) RedisClient() cache.RedisClient {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig())
	}

	return s.redisClient
}

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepositoryPg.NewRepository(s.DBClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) CacheRepository(_ context.Context) repository.CacheRepository {
	if s.cacheRepository == nil {
		s.cacheRepository = authRepositoryRedis.NewRepository(s.RedisClient())
	}

	return s.cacheRepository
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(s.AuthRepository(ctx), s.CacheRepository(ctx))
	}

	return s.authService
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}
