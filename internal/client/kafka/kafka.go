package kafka

import "context"

// Producer интерфейс кафка продьюсера
type Producer interface {
	Produce(ctx context.Context, data []byte) error
	Close() error
}
