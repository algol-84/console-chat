//
// Папка repository сожержит реализицию работы с репозиторием Postgres Auth

// Файл содержит только интерфейсы для работы с репо слоем. Здесь нет имплементации и поэтому нет зависимости от репо слоя.
// В репо может быть несколько репозиториев к разным таблицам или даже другим БД (не только постгрес), но все интерфейсы складываются сюда
// Вся имплементация уже располагается в подпапках (auth) с названием имени репозитория

package repository

import (
	"context"

	desc "github.com/algol-84/auth/pkg/user_v1"
)

type AuthRepository interface {
	Create(ctx context.Context, user *desc.User) (int64, error)
	Get(ctx context.Context, id int64) (*desc.UserInfo, error)
	Update(ctx context.Context, user *desc.UserUpdate) error
	Delete(ctx context.Context, id int64) error
}
