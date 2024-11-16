package auth

import (
	"context"

	desc "github.com/algol-84/auth/pkg/access_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Check ручка авторизации пользователя
func (i *Implementation) Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error) {
	err := i.accessService.Check(ctx, req.GetEndpointAddress())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "check error")
	}
	return &emptypb.Empty{}, nil
}
