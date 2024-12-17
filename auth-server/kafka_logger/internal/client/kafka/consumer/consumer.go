package consumer

import (
	"context"
	"log"
	"strings"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

type consumer struct {
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *GroupHandler
	topic                string
}

func NewConsumer(
	consumerGroup sarama.ConsumerGroup,
	consumerGroupHandler *GroupHandler,
	topic string,
) *consumer {
	return &consumer{
		consumerGroup:        consumerGroup,
		consumerGroupHandler: consumerGroupHandler,
		topic:                topic,
	}
}

func (c *consumer) Consume(ctx context.Context, handler Handler) error {
	c.consumerGroupHandler.msgHandler = handler

	return c.consume(ctx)
}

func (c *consumer) Close() error {
	return c.consumerGroup.Close()
}

func (c *consumer) consume(ctx context.Context) error {
	for {
		err := c.consumerGroup.Consume(ctx, strings.Split(c.topic, ","), c.consumerGroupHandler)
		if err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		log.Printf("rebalancing...\n")
	}
}
