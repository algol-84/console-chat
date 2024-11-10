package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/algol-84/auth/internal/model"
)

func (s *service) Create(ctx context.Context, user *model.User) (int64, error) {
	id, err := s.authRepository.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("repo error")
	}

	user.ID = id

	// Сериализовать структуру с пользователем в JSON
	data, err := json.Marshal(user)
	if err != nil {
		log.Printf("failed to marshal data: %v\n", err.Error())
		return 0, err
	}

	// Отправить JSON с юзером в кафку
	err = s.kafkaProducer.Produce(ctx, data)
	if err != nil {
		log.Printf("failed to produce log message to Kafka: %v\n", err.Error())
		return 0, err
	}

	return id, nil
}
