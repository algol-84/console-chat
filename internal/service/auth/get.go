package auth

import (
	"context"
	"fmt"

	"github.com/algol-84/auth/internal/model"
)

func (s *service) Get(ctx context.Context, id int64) (*model.User, error) {
	// Запрос юзера из кэша
	user, err := s.cacheRepository.Get(ctx, id)
	if err == model.ErrorUserNotFound {
		// Запрос юзера из базы
		user, err = s.authRepository.Get(ctx, id)
		if err != nil {
			return &model.User{}, model.ErrorUserNotFound
		}

		// Добавить юзера в кэш
		_, err = s.cacheRepository.Create(ctx, user)
		if err != nil {
			return &model.User{}, fmt.Errorf("cache error")
		}

		return user, nil
	}

	return user, nil
}
