package config

import (
	"errors"
	"os"
	"strings"
)

const (
	brokersEnvName = "KAFKA_BROKERS"
	groupIDEnvName = "KAFKA_GROUP_ID"
	topicEnvName   = "KAFKA_TOPIC"
)

type kafkaProducerConfig struct {
	brokers []string
	groupID string
	topic   string
}

// NewKafkaProducerConfig читает настройки кафки из файла конфига
func NewKafkaProducerConfig() (*kafkaProducerConfig, error) {
	brokersStr := os.Getenv(brokersEnvName)
	if len(brokersStr) == 0 {
		return nil, errors.New("kafka brokers address not found")
	}

	brokers := strings.Split(brokersStr, ",")

	groupID := os.Getenv(groupIDEnvName)
	if len(groupID) == 0 {
		return nil, errors.New("kafka group id not found")
	}

	topic := os.Getenv(topicEnvName)
	if len(topic) == 0 {
		return nil, errors.New("kafka topic name not found")
	}

	return &kafkaProducerConfig{
		brokers: brokers,
		groupID: groupID,
		topic:   topic,
	}, nil
}

func (cfg *kafkaProducerConfig) Brokers() []string {
	return cfg.brokers
}

func (cfg *kafkaProducerConfig) GroupID() string {
	return cfg.groupID
}

func (cfg *kafkaProducerConfig) Topic() string {
	return cfg.topic
}
