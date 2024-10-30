package auth

import (
	"context"
	"errors"

	"github.com/algol-84/auth/internal/model"
)

func (s *service) Update(ctx context.Context, user *model.UserUpdate) error {
	// Проверить что запрос не пустой
	if !user.Name.Valid && !user.Email.Valid && !user.Role.Valid {
		return errors.New("update request is empty")
	}

	err := s.authRepository.Update(ctx, user)
	if err != nil {
		return err
	}

	// Удалить юзера из кэша
	err = s.cacheRepository.Delete(ctx, user.ID)
	if err != nil {
		return err
	}

	return nil
}
