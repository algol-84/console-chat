package app

import (
	"context"
	"log"

	redigo "github.com/gomodule/redigo/redis"

	accessApi "github.com/algol-84/auth/internal/api/access"
	authApi "github.com/algol-84/auth/internal/api/auth"
	userApi "github.com/algol-84/auth/internal/api/user"
	"github.com/algol-84/auth/internal/client/cache"
	"github.com/algol-84/auth/internal/client/cache/redis"
	"github.com/algol-84/auth/internal/client/kafka"
	"github.com/algol-84/auth/internal/client/kafka/producer"
	"github.com/algol-84/auth/internal/config"
	"github.com/algol-84/auth/internal/repository"
	authRepositoryPg "github.com/algol-84/auth/internal/repository/auth/pg"
	authRepositoryRedis "github.com/algol-84/auth/internal/repository/auth/redis"
	"github.com/algol-84/auth/internal/service"
	accessService "github.com/algol-84/auth/internal/service/access"
	authService "github.com/algol-84/auth/internal/service/auth"
	userService "github.com/algol-84/auth/internal/service/user"
	closer "github.com/algol-84/platform_common/pkg/closer"
	db "github.com/algol-84/platform_common/pkg/db"
	pg "github.com/algol-84/platform_common/pkg/db/pg"
)

// serviceProvider хранит все объекты приложения, как интерфейсы или ссылки на структуры
type serviceProvider struct {
	pgConfig            config.PGConfig
	grpcConfig          config.GRPCConfig
	redisConfig         config.RedisConfig
	kafkaProducerConfig config.KafkaProducerConfig

	kafkaProducer   kafka.Producer
	dbClient        db.Client
	redisPool       *redigo.Pool
	redisClient     cache.RedisClient
	authRepository  repository.AuthRepository
	cacheRepository repository.CacheRepository

	userService   service.UserService
	authService   service.AuthService
	accessService service.AccessService

	userImpl   *userApi.Implementation
	authImpl   *authApi.Implementation
	accessImpl *accessApi.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// Определяются функции инициализации всех объектов

func (s *serviceProvider) TokenConfig() {
	
}

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

// KafkaProducerConfig инициализирует считывание настроек Кафки из файла конфига
func (s *serviceProvider) KafkaProducerConfig() config.KafkaProducerConfig {
	if s.kafkaProducerConfig == nil {
		cfg, err := config.NewKafkaProducerConfig()
		if err != nil {
			log.Fatalf("failed to get kafka consumer config: %s", err.Error())
		}

		s.kafkaProducerConfig = cfg
	}

	return s.kafkaProducerConfig
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

func (s *serviceProvider) KafkaProducer() kafka.Producer {
	if s.kafkaProducer == nil {
		s.kafkaProducer = producer.NewProducer(s.KafkaProducerConfig())
	}

	return s.kafkaProducer
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

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.AuthRepository(ctx), s.CacheRepository(ctx), s.KafkaProducer())
	}

	return s.userService
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(s.AuthRepository(ctx))
	}

	return s.authService
}

func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessService.NewService(s.AuthRepository(ctx))
	}

	return s.accessService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *userApi.Implementation {
	if s.userImpl == nil {
		s.userImpl = userApi.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *authApi.Implementation {
	if s.authImpl == nil {
		s.authImpl = authApi.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}

func (s *serviceProvider) AccessImpl(ctx context.Context) *accessApi.Implementation {
	if s.accessImpl == nil {
		s.accessImpl = accessApi.NewImplementation(s.AccessService(ctx))
	}

	return s.accessImpl
}
