package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/algol-84/auth/internal/converter"
	desc "github.com/algol-84/auth/pkg/user_v1"
)

// func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
// 	id, err := i.noteService.Create(ctx, converter.ToNoteInfoFromDesc(req.GetInfo()))
// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Printf("inserted note with id: %d", id)

// 	return &desc.CreateResponse{
// 		Id: id,
// 	}, nil
// }

// Create обрабатывает GRPC запросы на создание нового юзера
func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	userID, err := i.authService.Create(ctx, converter.FromUserToService(req.User))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "user creation in DB returned with error: %s", err)
	}

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}
