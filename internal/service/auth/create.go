package auth

import (
	"context"
	"fmt"

	"github.com/algol-84/auth/internal/model"
)

func (s *service) Create(ctx context.Context, user *model.User) (int64, error) {
	id, err := s.authRepository.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("repo error")
	}

	return id, nil
}
