package model

import (
	"time"
)

// Chat представляет структуру списка чатов сервисного слоя
type Chat struct {
	ID        int64
	Usernames []string
	CreatedAt time.Time
}
