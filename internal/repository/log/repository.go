package log

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/algol-84/chat-server/internal/client/db"
	"github.com/algol-84/chat-server/internal/repository"
	// "github.com/algol-84/chat-server/internal/repository/chat/converter"
	//model "github.com/algol-84/chat-server/internal/model"
	// modelRepo "github.com/algol-84/chat-server/internal/repository/chat/model"
)

// Представление БД сервиса Log
const chatTable = "chat_log" // таблица хранит лог

const (
	fieldChatID    = "chat_id"
	fieldInfo      = "info"
	fieldCreatedAt = "created_at"
)

type repo struct {
	db db.Client
}

// NewRepository конструктор
func NewRepository(db db.Client) repository.LogRepository {
	return &repo{db: db}
}

// Create добавляет запись в лог с id чата и информацией о действии
func (r *repo) Create(ctx context.Context, id int64, info string) (int64, error) {
	var logID int64
	// Собрать запрос на вставку записи в таблицу
	builderQuery := sq.Insert(chatTable).
		PlaceholderFormat(sq.Dollar).
		Columns(fieldChatID, fieldInfo).
		Values(id, info)

	query, args, err := builderQuery.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "log_repository.Create",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return 0, err
	}

	return logID, nil
}
