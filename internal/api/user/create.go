package auth

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/algol-84/auth/internal/converter"
	desc "github.com/algol-84/auth/pkg/user_v1"

	"github.com/algol-84/auth/internal/logger"
)

// Create обрабатывает GRPC запросы на создание нового юзера
func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	logger.Info("Creare user...", zap.String("user", req.User.Name))
	userID, err := i.userService.Create(ctx, converter.FromUserToService(req.User))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "user creation in DB returned with error")
	}

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}
