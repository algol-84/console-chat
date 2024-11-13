package auth

import (
	"context"

	desc "github.com/algol-84/auth/pkg/auth_v1"
)

// GetAccessToken ручка получения акцесс токена
func (i *Implementation) GetAccessToken(_ context.Context, _ *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {

	return &desc.GetAccessTokenResponse{
		AccessToken: "access token",
	}, nil
}
