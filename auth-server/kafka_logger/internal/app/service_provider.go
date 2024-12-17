package app

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/algol-84/kafka_logger/internal/client/kafka"
	kafkaConsumer "github.com/algol-84/kafka_logger/internal/client/kafka/consumer"
	"github.com/algol-84/kafka_logger/internal/config"
	"github.com/algol-84/kafka_logger/internal/config/env"
	closer "github.com/algol-84/platform_common/pkg/closer"

	"github.com/algol-84/kafka_logger/internal/service"
	noteSaverConsumer "github.com/algol-84/kafka_logger/internal/service/consumer/note_saver"
)

type serviceProvider struct {
	kafkaConsumerConfig config.KafkaConsumerConfig

	noteSaverConsumer service.ConsumerService

	consumer             kafka.Consumer
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *kafkaConsumer.GroupHandler
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) KafkaConsumerConfig() config.KafkaConsumerConfig {
	if s.kafkaConsumerConfig == nil {
		cfg, err := env.NewKafkaConsumerConfig()
		if err != nil {
			log.Fatalf("failed to get kafka consumer config: %s", err.Error())
		}

		s.kafkaConsumerConfig = cfg
	}

	return s.kafkaConsumerConfig
}

func (s *serviceProvider) NoteSaverConsumer(ctx context.Context) service.ConsumerService {
	if s.noteSaverConsumer == nil {
		s.noteSaverConsumer = noteSaverConsumer.NewService(
			s.Consumer(),
		)
	}

	return s.noteSaverConsumer
}

func (s *serviceProvider) Consumer() kafka.Consumer {
	if s.consumer == nil {
		s.consumer = kafkaConsumer.NewConsumer(
			s.ConsumerGroup(),
			s.ConsumerGroupHandler(),
			s.kafkaConsumerConfig.Topic(),
		)
		closer.Add(s.consumer.Close)
	}

	return s.consumer
}

func (s *serviceProvider) ConsumerGroup() sarama.ConsumerGroup {
	if s.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			s.KafkaConsumerConfig().Brokers(),
			s.KafkaConsumerConfig().GroupID(),
			s.KafkaConsumerConfig().Config(),
		)
		if err != nil {
			log.Fatalf("failed to create consumer group: %v", err)
		}

		s.consumerGroup = consumerGroup
	}

	return s.consumerGroup
}

func (s *serviceProvider) ConsumerGroupHandler() *kafkaConsumer.GroupHandler {
	if s.consumerGroupHandler == nil {
		s.consumerGroupHandler = kafkaConsumer.NewGroupHandler()
	}

	return s.consumerGroupHandler
}
