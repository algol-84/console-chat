package auth

import (
	"context"

	desc "github.com/algol-84/auth/pkg/auth_v1"
)

// GetRefreshToken ручка получения рефреш токена
func (i *Implementation) GetRefreshToken(_ context.Context, _ *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {

	return &desc.GetRefreshTokenResponse{
		RefreshToken: "new token",
	}, nil
}
