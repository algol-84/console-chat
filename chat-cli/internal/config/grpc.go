package config

import (
	"errors"
	"net"
	"os"
)

const (
	authHostEnvName = "AUTH_HOST"
	authPortEnvName = "AUTH_PORT"
	chatHostEnvName = "CHAT_HOST"
	chatPortEnvName = "CHAT_PORT"
)

// GRPCConfig содержит настройки для gRPC сервера
type GRPCConfig interface {
	AuthServiceAddress() string
	ChatServiceAddress() string
}

type grpcConfig struct {
	authServiceHost string
	authServicePort string
	chatServiceHost string
	chatServicePort string
}

// NewGRPCConfig считывает аргументы из конфига и заполняет grpcConfig
func NewGRPCConfig() (GRPCConfig, error) {
	authHost := os.Getenv(authHostEnvName)
	if len(authHost) == 0 {
		return nil, errors.New("auth service host not found")
	}

	authPort := os.Getenv(authPortEnvName)
	if len(authPort) == 0 {
		return nil, errors.New("auth service port not found")
	}

	chatHost := os.Getenv(chatHostEnvName)
	if len(chatHost) == 0 {
		return nil, errors.New("chat service host not found")
	}

	chatPort := os.Getenv(chatPortEnvName)
	if len(chatPort) == 0 {
		return nil, errors.New("chat service port not found")
	}

	return &grpcConfig{
		authServiceHost: authHost,
		authServicePort: authPort,
		chatServiceHost: chatHost,
		chatServicePort: chatPort,
	}, nil
}

func (cfg *grpcConfig) AuthServiceAddress() string {
	return net.JoinHostPort(cfg.authServiceHost, cfg.authServicePort)
}

func (cfg *grpcConfig) ChatServiceAddress() string {
	return net.JoinHostPort(cfg.chatServiceHost, cfg.chatServicePort)
}
