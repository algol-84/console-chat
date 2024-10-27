package app

import (
	"context"
	"log"

	chat "github.com/algol-84/chat-server/internal/api/chat"
	"github.com/algol-84/chat-server/internal/client/db"
	"github.com/algol-84/chat-server/internal/client/db/pg"
	"github.com/algol-84/chat-server/internal/client/db/transaction"
	"github.com/algol-84/chat-server/internal/closer"
	"github.com/algol-84/chat-server/internal/config"
	"github.com/algol-84/chat-server/internal/repository"
	chatRepository "github.com/algol-84/chat-server/internal/repository/chat"
	logRepository "github.com/algol-84/chat-server/internal/repository/log"
	"github.com/algol-84/chat-server/internal/service"
	chatService "github.com/algol-84/chat-server/internal/service/chat"
)

// serviceProvider хранит все объекты приложения, как интерфейсы или ссылки на структуры
type serviceProvider struct {
	pgConfig       config.PGConfig
	grpcConfig     config.GRPCConfig
	dbClient       db.Client
	txManager      db.TxManager
	chatRepository repository.ChatRepository
	logRepository  repository.LogRepository
	chatService    service.ChatService
	chatImpl       *chat.Implementation
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

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewRepository(s.DBClient(ctx))
	}

	return s.chatRepository
}

func (s *serviceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepository.NewRepository(s.DBClient(ctx))
	}

	return s.logRepository
}

func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(s.ChatRepository(ctx), s.LogRepository(ctx), s.TxManager(ctx))
	}

	return s.chatService
}

func (s *serviceProvider) ChatImpl(ctx context.Context) *chat.Implementation {
	if s.chatImpl == nil {
		s.chatImpl = chat.NewImplementation(s.ChatService(ctx))
	}

	return s.chatImpl
}
