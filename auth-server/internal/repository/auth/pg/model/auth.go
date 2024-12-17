//
// Файл содержит структуры моделей представлений таблиц БД
// Это внутренний тип данных внутри репо слоя

package model

import (
	"database/sql"
	"time"
)

// User представляет собой модель репо слоя
type User struct {
	ID        int64        `db:"id"`
	Name      string       `db:"name"`
	Password  string       `db:"password"`
	Email     string       `db:"email"`
	Role      string       `db:"role"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
