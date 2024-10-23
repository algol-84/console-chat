package chat

import (
	"github.com/algol-84/chat-server/internal/client/db"
	"github.com/algol-84/chat-server/internal/repository"
	def "github.com/algol-84/chat-server/internal/service"
)

type service struct {
	chatRepository repository.ChatRepository
	logRepository  repository.LogRepository
	txManager      db.TxManager
}

// NewService конструктор сервисного слоя
func NewService(chatRepository repository.ChatRepository,
	logRepository repository.LogRepository,
	txManager db.TxManager) def.ChatService {
	return &service{
		chatRepository: chatRepository,
		logRepository:  logRepository,
		txManager:      txManager,
	}
}
