package chat

import (
	"github.com/algol-84/chat-server/internal/repository"
	def "github.com/algol-84/chat-server/internal/service"
)

type service struct {
	chatRepository repository.ChatRepository
}

// NewService конструктор сервисного слоя
func NewService(chatRepository repository.ChatRepository) def.ChatService {
	return &service{
		chatRepository: chatRepository,
	}
}

func NewMockService(deps ...interface{}) def.ChatService {
	srv := service{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.ChatRepository:
			srv.chatRepository = s
		}
	}

	return &srv
}
