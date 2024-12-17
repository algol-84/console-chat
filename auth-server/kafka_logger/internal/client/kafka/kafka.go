package kafka

import (
	"context"

	"github.com/algol-84/kafka_logger/internal/client/kafka/consumer"
)

type Consumer interface {
	Consume(ctx context.Context, handler consumer.Handler) (err error)
	Close() error
}
