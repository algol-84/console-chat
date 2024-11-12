package note_saver

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"

	"github.com/algol-84/kafka_logger/internal/model"
)

func (s *service) AuthSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	user := &model.User{}
	err := json.Unmarshal(msg.Value, user)
	if err != nil {
		return err
	}

	log.Printf("consumed from kafka: %v", user)

	return nil
}
