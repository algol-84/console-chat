package auth

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/algol-84/chat-server/pkg/chat_v1"
	"github.com/google/uuid"
)

// Create обрабатывает GRPC запросы на создание нового чата
func (i *Implementation) Create(ctx context.Context, _ *emptypb.Empty) (*desc.CreateResponse, error) {
	// span, _ := opentracing.StartSpanFromContext(ctx, "create api")
	// defer span.Finish()

	chatID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	log.Println(chatID.ID(), chatID.String())

	i.channels[chatID.String()] = make(chan *desc.Message, 100)

	return &desc.CreateResponse{
		Id: 1, //chatID.String(),
	}, nil
}
