package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

const (
	brokersEnvName = "KAFKA_BROKERS"
	groupIDEnvName = "KAFKA_GROUP_ID"
	topicEnvName   = "KAFKA_TOPIC"
	retryMaxName   = "KAFKA_RETRY_MAX"
)

type kafkaProducerConfig struct {
	brokers  []string
	groupID  string
	topic    string
	retryMax int
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

	retryMax, err := strconv.Atoi(os.Getenv(retryMaxName))
	if err != nil {
		return nil, errors.New("kafka retryMax name not found")
	}

	return &kafkaProducerConfig{
		brokers:  brokers,
		groupID:  groupID,
		topic:    topic,
		retryMax: retryMax,
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

func (cfg *kafkaProducerConfig) RetryMax() int {
	return cfg.retryMax
}
