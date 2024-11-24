package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	conv "github.com/algol-84/chat-server/internal/converter"
	desc "github.com/algol-84/chat-server/pkg/chat_v1"
	"github.com/opentracing/opentracing-go"
)

// Create обрабатывает GRPC запросы на создание нового юзера
func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "create api")
	defer span.Finish()

	userID, err := i.chatService.Create(ctx, conv.FromChatToService(req.Chat))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "user creation in DB returned with error")
	}

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}
