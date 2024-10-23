package converter

import (
	model "github.com/algol-84/chat-server/internal/model"
	modelRepo "github.com/algol-84/chat-server/internal/repository/chat/model"
)

// ToChatFromRepo конвертирует из модели репо слоя в модель сервисного слоя
func ToChatFromRepo(chat *modelRepo.Chat) *model.Chat {
	return &model.Chat{
		ID:        chat.ID,
		CreatedAt: chat.CreatedAt,
	}
}
