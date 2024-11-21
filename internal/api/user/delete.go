package auth

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/algol-84/auth/internal/logger"
	desc "github.com/algol-84/auth/pkg/user_v1"
)

// Delete обрабатывает GRPC запросы на удаление пользователя
func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	logger.Info("Delete user...", zap.Int64("user id:", req.Id))
	err := i.userService.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "removing user from the DB returned with an error")
	}

	return &emptypb.Empty{}, nil
}
