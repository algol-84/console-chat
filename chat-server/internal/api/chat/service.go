package chat

import (
	"sync"

	"github.com/algol-84/chat-server/internal/service"
	desc "github.com/algol-84/chat-server/pkg/chat_v1"
)

// Chat определяет мапу стримов для каждого чата и мьютекс для раздельного доступа к каждому стриму
type Chat struct {
	streams map[string]desc.ChatV1_ConnectChatServer
	m       sync.RWMutex
}

// Implementation содержит все объекты сервера
type Implementation struct {
	desc.UnimplementedChatV1Server
	chatService service.ChatService

	chats  map[string]*Chat
	mxChat sync.RWMutex

	channels  map[string]chan *desc.Message
	mxChannel sync.RWMutex
}

// NewImplementation конструктор API слоя
func NewImplementation(chatService service.ChatService) *Implementation {
	return &Implementation{
		chatService: chatService,
		chats:       make(map[string]*Chat),
		channels:    make(map[string]chan *desc.Message),
	}
}
