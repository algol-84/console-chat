package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/algol-84/auth/internal/logger"
	"github.com/algol-84/auth/internal/model"
	"go.uber.org/zap"
)

func (s *service) Create(ctx context.Context, user *model.User) (int64, error) {
	id, err := s.authRepository.Create(ctx, user)
	if err != nil {
		logger.Error("failed to create user", zap.Int64("id", user.ID))
		return 0, fmt.Errorf("repo error")
	}

	user.ID = id

	// Сериализовать структуру с пользователем в JSON
	data, err := json.Marshal(user)
	if err != nil {
		logger.Error("failed to marshal data", zap.String("error", err.Error()))
		return 0, err
	}

	// Отправить JSON с юзером в кафку
	err = s.kafkaProducer.Produce(ctx, data)
	if err != nil {
		logger.Error("failed to produce log message to Kafka", zap.String("error", err.Error()))
		return 0, err
	}

	return id, nil
}
