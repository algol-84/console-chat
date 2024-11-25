package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/algol-84/auth/internal/converter"
	"github.com/algol-84/auth/internal/logger"
	desc "github.com/algol-84/auth/pkg/user_v1"
)

// Create обрабатывает GRPC запросы на создание нового юзера
func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	logger.Info("Creare user...")
	userID, err := i.userService.Create(ctx, converter.FromUserToService(req.User))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "user creation in DB returned with error")
	}

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}
