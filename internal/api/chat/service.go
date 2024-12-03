package auth

import (
	"github.com/algol-84/chat-server/internal/service"
	desc "github.com/algol-84/chat-server/pkg/chat_v1"
)

// Implementation содержит все объекты сервера
type Implementation struct {
	desc.UnimplementedChatV1Server
	chatService service.ChatService
	channels    map[string]chan *desc.Message
}

// NewImplementation конструктор API слоя
func NewImplementation(chatService service.ChatService) *Implementation {
	return &Implementation{
		chatService: chatService,
		channels: make(map[string]chan *desc.Message),
	}
}
