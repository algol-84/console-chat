package auth

import (
	"context"

	desc "github.com/algol-84/auth/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login ручка получения первичного рефреш токена
func (i *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {

	refreshToken, err := i.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &desc.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
