package auth

import (
	"context"

	desc "github.com/algol-84/auth/pkg/auth_v1"
)

// Login ручка получения первичного рефреш токена
func (i *Implementation) Login(_ context.Context, _ *desc.LoginRequest) (*desc.LoginResponse, error) {

	return &desc.LoginResponse{
		RefreshToken: "token",
	}, nil
}
