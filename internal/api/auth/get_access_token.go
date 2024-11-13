package auth

import (
	"context"

	desc "github.com/algol-84/auth/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAccessToken ручка получения акцесс токена
func (i *Implementation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	accessToken, err := i.authService.GetAccessToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "invalid access token")
	}

	return &desc.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
