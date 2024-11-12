package config

import (
	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

type KafkaConsumerConfig interface {
	Brokers() []string
	GroupID() string
	Topic() string
	Config() *sarama.Config
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
