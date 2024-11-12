package note_saver

import (
	"context"
	"log"

	"github.com/algol-84/kafka_logger/internal/client/kafka"
	def "github.com/algol-84/kafka_logger/internal/service"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	consumer kafka.Consumer
}

func NewService(
	consumer kafka.Consumer,
) *service {
	return &service{
		consumer: consumer,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	log.Printf("kafka consumer is running")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-s.run(ctx):
			if err != nil {
				return err
			}
		}
	}
}

func (s *service) run(ctx context.Context) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		errChan <- s.consumer.Consume(ctx, s.AuthSaveHandler)
	}()

	return errChan
}
