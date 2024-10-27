// Package converter Конвертер типов protobuf <-> service model
package converter

import (
	"github.com/algol-84/chat-server/internal/model"
	desc "github.com/algol-84/chat-server/pkg/chat_v1"
)

// FromChatToService конвертирует тип protobuf User в модель сервисного слоя
func FromChatToService(chat *desc.Chat) *model.Chat {
	return &model.Chat{
		Usernames: chat.Usernames,
	}
}
