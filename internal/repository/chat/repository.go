package chat

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/algol-84/chat-server/internal/repository"
	db "github.com/algol-84/platform_common/pkg/db"

	model "github.com/algol-84/chat-server/internal/model"
)

// Представление БД сервиса Chat-server
const chatTable = "chat_room" // таблица хранит список созданных чатов
const userTable = "chat_user" // таблица хранит сопоставление юзеров и чатов

const (
	fieldChatID    = "chat_id"
	fieldCreatedAt = "created_at"
	fieldUserName  = "user_name"
)

type repo struct {
	db db.Client
}

// NewRepository конструктор
func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

// Create создает новый чат в БД
// Функция возвращает присвоенный в БД ID чата или ошибку записи
func (r *repo) Create(ctx context.Context, chat *model.Chat) (int64, error) {
	for _, username := range chat.Usernames {
		if username == "" {
			return 0, errors.New("username must not be empty")
		}
	}

	var chatID int64
	// Собрать запрос на вставку записи в таблицу
	builderQuery := sq.Insert(chatTable).
		PlaceholderFormat(sq.Dollar).
		Columns(fieldCreatedAt).
		Values(time.Now()).
		Suffix("RETURNING " + fieldChatID)

	query, args, err := builderQuery.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "chat_repository.Create",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID)
	if err != nil {
		return 0, err
	}

	builderQuery = sq.Insert(userTable).
		PlaceholderFormat(sq.Dollar).
		Columns(fieldChatID, fieldUserName)

	for _, username := range chat.Usernames {
		builderQuery = builderQuery.Values(chatID, username)
	}
	query, args, err = builderQuery.ToSql()
	if err != nil {
		return 0, err
	}

	q = db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

// Delete удаляет чат и юзеров из этого чата в связанной таблице chat_user
func (r *repo) Delete(ctx context.Context, chatID int64) error {
	builderQuery := sq.Delete(chatTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{fieldChatID: chatID})

	query, args, err := builderQuery.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID=%d not found", chatID)
	}

	return nil
}
