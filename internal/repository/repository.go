//
// Папка repository сожержит реализицию работы с репозиторием Postgres Auth

// Файл содержит только интерфейсы для работы с репо слоем. Здесь нет имплементации и поэтому нет зависимости от репо слоя.
// В репо может быть несколько репозиториев к разным таблицам или даже другим БД (не только постгрес), но все интерфейсы складываются сюда
// Вся имплементация уже располагается в подпапках (auth) с названием имени репозитория

package repository

import (
	"context"

	model "github.com/algol-84/auth/internal/model"
)

// AuthRepository интерфейс реализует репо слой User репозитория
type AuthRepository interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error
	Find(ctx context.Context, username string, password string) (*model.User, error)
}

// CacheRepository интерфейс реализует репо слой кэша
type CacheRepository interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Delete(ctx context.Context, id int64) error
}

// AccessRepository интерфейс репо слоя Access репозитория
type AccessRepository interface {
	Get(ctx context.Context) (map[string]string, error)
}
