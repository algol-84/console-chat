//
// Папка repository сожержит реализицию работы с репозиторием Postgres

// Файл содержит только интерфейсы для работы с репо слоем. Здесь нет имплементации и поэтому нет зависимости от репо слоя.
// В репо может быть несколько репозиториев к разным таблицам или даже другим БД (не только постгрес), но все интерфейсы складываются сюда
// Вся имплементация уже располагается в подпапках с названием имени репозитория

package repository

import (
	"context"

	"github.com/algol-84/chat-server/internal/model"
)

// ChatRepository интерфейс реализует репо слой
type ChatRepository interface {
	Create(ctx context.Context, chat *model.Chat) (int64, error)
	Delete(ctx context.Context, id int64) error
}
