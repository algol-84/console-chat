package chat

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/algol-84/chat-server/pkg/chat_v1"
	"github.com/google/uuid"
)

// CreateChat обрабатывает GRPC запросы на создание нового чата
func (i *Implementation) CreateChat(_ context.Context, _ *emptypb.Empty) (*desc.CreateChatResponse, error) {
	chatID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	i.channels[chatID.String()] = make(chan *desc.Message, 100)
	log.Printf("chat created with id: %s", chatID)

	return &desc.CreateChatResponse{
		Id: chatID.String(),
	}, nil
}
