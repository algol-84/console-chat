package chat

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/algol-84/chat-server/pkg/chat_v1"
)

// DeleteChat обрабатывает GRPC запросы на удаление пользователя
func (i *Implementation) DeleteChat(_ context.Context, req *desc.DeleteChatRequest) (*emptypb.Empty, error) {
	// Delete channel from map
	delete(i.channels, req.Id)
	log.Printf("chat with id: %s is deleted", req.Id)

	return &emptypb.Empty{}, nil
}
