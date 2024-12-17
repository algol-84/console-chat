package producer

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/algol-84/auth/internal/config"
)

type producer struct {
	producer sarama.SyncProducer
	config   config.KafkaProducerConfig
}

// NewProducer конструктор кафка продьюсера
func NewProducer(config config.KafkaProducerConfig) *producer {
	p, err := newSyncProducer(config.Brokers(), config.RetryMax())
	if err != nil {
		log.Printf("failed to start producer: %v\n", err.Error())
	}

	return &producer{
		producer: p,
		config:   config,
	}
}

func newSyncProducer(brokerList []string, retryMax int) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = retryMax
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

// Produce отправляет сообщение в кафку
func (c *producer) Produce(_ context.Context, data []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: c.config.Topic(),
		Value: sarama.StringEncoder(data),
	}

	partition, offset, err := c.producer.SendMessage(msg)
	if err != nil {
		log.Printf("failed to send message in Kafka: %v\n", err.Error())
		return err
	}

	log.Printf("message sent to kafka topic %s in partition %d with offset %d\n", msg.Topic, partition, offset)

	return nil
}

func (c *producer) Close() error {
	log.Printf("close kafka")
	if err := c.producer.Close(); err != nil {
		log.Printf("failed to close producer: %v\n", err.Error())
	}
	return nil
}
