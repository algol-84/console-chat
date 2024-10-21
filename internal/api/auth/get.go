package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/algol-84/auth/internal/converter"
	desc "github.com/algol-84/auth/pkg/user_v1"
)

// func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
// 	noteObj, err := i.noteService.Get(ctx, req.GetId())
// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", noteObj.ID, noteObj.Info.Title, noteObj.Info.Content, noteObj.CreatedAt, noteObj.UpdatedAt)

// 	return &desc.GetResponse{
// 		Note: converter.ToNoteFromService(noteObj),
// 	}, nil
// }

// Get обрабатывает GRPC запросы на получение данных пользователя
func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := i.authService.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "the request for user data in the DB returned with error: %s", err)
	}

	return &desc.GetResponse{
		UserInfo: converter.ToUserInfoFromService(user),
	}, nil
}
