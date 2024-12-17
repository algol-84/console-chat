package env

import (
	"errors"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

const (
	brokersEnvName = "KAFKA_BROKERS"
	groupIDEnvName = "KAFKA_GROUP_ID"
	topicEnvName   = "KAFKA_TOPIC"
)

type kafkaConsumerConfig struct {
	brokers []string
	groupID string
	topic   string
}

func NewKafkaConsumerConfig() (*kafkaConsumerConfig, error) {
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

	return &kafkaConsumerConfig{
		brokers: brokers,
		groupID: groupID,
		topic:   topic,
	}, nil
}

func (cfg *kafkaConsumerConfig) Brokers() []string {
	return cfg.brokers
}

func (cfg *kafkaConsumerConfig) GroupID() string {
	return cfg.groupID
}

func (cfg *kafkaConsumerConfig) Topic() string {
	return cfg.topic
}

// Config возвращает конфигурацию для sarama consumer
func (cfg *kafkaConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
