package auth

import (
	"context"

	desc "github.com/algol-84/auth/pkg/access_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Check ручка авторизации пользователя
func (i *Implementation) Check(_ context.Context, _ *desc.CheckRequest) (*emptypb.Empty, error) {

	return &emptypb.Empty{}, nil
}
