package config

import (
	"time"

	"github.com/joho/godotenv"
)

// Load загружает конфиг из файла конфигурации
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

// RedisConfig интерфейс конфига редиса
type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

// KafkaProducerConfig интерфейс конфига кафки
type KafkaProducerConfig interface {
	Brokers() []string
	GroupID() string
	Topic() string
	RetryMax() int
}
