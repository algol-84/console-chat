// Package model содержит структуры моделей представлений таблиц БД
// Это внутренний тип данных внутри репо слоя
package model

import (
	"time"
)

// Chat представляет собой модель репо слоя
type Chat struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}
