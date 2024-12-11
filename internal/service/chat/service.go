package chat

import (
	def "github.com/algol-84/chat-server/internal/service"
)

type service struct {
}

// NewService конструктор сервисного слоя
func NewService() def.ChatService {
	return &service{}
}
