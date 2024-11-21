package auth

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/algol-84/auth/internal/converter"
	"github.com/algol-84/auth/internal/logger"
	desc "github.com/algol-84/auth/pkg/user_v1"
)

// Update обрабатывает GRPC запросы на обновление данных пользователя
func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	logger.Info("Update user...", zap.Int64("user id:", req.UserUpdate.Id))
	err := i.userService.Update(ctx, converter.ToUserUpdateFromDesc(req.UserUpdate))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "updating user in the DB returned with an error")
	}

	return &emptypb.Empty{}, nil
}
