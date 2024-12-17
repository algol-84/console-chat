package app

import (
	"context"
	"log"

	chat "github.com/algol-84/chat-server/internal/api/chat"
	chatService "github.com/algol-84/chat-server/internal/service/chat"

	"github.com/algol-84/chat-server/internal/config"
	"github.com/algol-84/chat-server/internal/service"
)

// serviceProvider хранит все объекты приложения, как интерфейсы или ссылки на структуры
type serviceProvider struct {
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig
	chatService   service.ChatService
	chatImpl      *chat.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// Определяются функции инициализации всех объектов

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

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) ChatService(_ context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService()
	}

	return s.chatService
}

func (s *serviceProvider) ChatImpl(ctx context.Context) *chat.Implementation {
	if s.chatImpl == nil {
		s.chatImpl = chat.NewImplementation(s.ChatService(ctx))
	}

	return s.chatImpl
}
