package pg_chat

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Представление БД сервиса Chat-server
const chatTable = "chat_room" // таблица хранит список созданных чатов
const userTable = "chat_user" // таблица хранит сопоставление юзеров и чатов

const (
	fieldChatID    = "chat_id"
	fieldCreatedAt = "created_at"
	fieldUserName  = "user_name"
)

// Указатель на пул соединений к БД
var pool *pgxpool.Pool

// User предсталяет пользователя в БД
// Содержит основные данные пользователя, такие как имя, электронная почта,
// хеш пароля, ID роли, а также время создания и обновления записи.
type User struct {
	ID        int64
	Name      string
	Email     string // Может принимать NULL
	Password  string
	Role      string
	CreatedAt time.Time // Заполняется в момент создания юзера
	UpdatedAt time.Time // Заполняется при апдейте, может принимать NULL
}

// Connect создает пул подключений к БД Auth
func Connect(ctx context.Context, connString string) error {
	pgxpool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return err
	}
	pool = pgxpool
	return nil
}

// Close закрывает пул соединений с БД
func Close() {
	pool.Close()
}

// CreateChat создает новый чат в БД
// Функция возвращает присвоенный в БД ID чата или ошибку записи
func CreateChat(ctx context.Context, usernames []string) (int64, error) {
	for _, username := range usernames {
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

	err = pool.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		return 0, err
	}

	builderQuery = sq.Insert(userTable).
		PlaceholderFormat(sq.Dollar).
		Columns(fieldChatID, fieldUserName)

	for _, username := range usernames {
		builderQuery = builderQuery.Values(chatID, username)

		log.Println(query, args)
	}
	query, args, err = builderQuery.ToSql()
	if err != nil {
		return 0, err
	}
	_, err = pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

// DeleteChat удаляет чат
func DeleteChat(ctx context.Context, chatID int64) error {
	builderQuery := sq.Delete(chatTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{fieldChatID: chatID})

	query, args, err := builderQuery.ToSql()
	if err != nil {
		return err
	}

	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID=%d not found", chatID)
	}

	return nil
}
