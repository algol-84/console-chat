package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/algol-84/auth/pkg/user_v1"
)

// Delete обрабатывает GRPC запросы на удаление пользователя
func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.authService.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "removing user from the DB returned with an error")
	}

	return &emptypb.Empty{}, nil
}
