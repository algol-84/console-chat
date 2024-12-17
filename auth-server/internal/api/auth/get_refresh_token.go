package auth

import (
	"context"

	desc "github.com/algol-84/auth/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRefreshToken ручка получения рефреш токена
func (i *Implementation) GetRefreshToken(ctx context.Context, req *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {
	refreshToken, err := i.authService.GetRefreshToken(ctx, req.OldRefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "invalid refresh token")
	}
	return &desc.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}
