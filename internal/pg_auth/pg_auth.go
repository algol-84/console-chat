package pg_auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"

	desc "github.com/algol-84/auth/pkg/user_v1"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type stringValue *wrapperspb.StringValue

// Представление БД сервиса Auth
const table = "chat_user"
const (
	fieldID        = "id"
	fieldName      = "name"
	fieldEmail     = "email"
	fieldRole      = "role"
	fieldPassword  = "password"
	fieldCreatedAt = "created_at"
	fieldUpdatedAt = "updated_at"
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
func Connect(ctx context.Context, connString string) (err error) {
	pool, err = pgxpool.Connect(ctx, connString)
	if err != nil {
		return err
	}
	return nil
}

// Close закрывает пул соединений с БД
func Close() {
	pool.Close()
}

// CreateUser создает нового юзера в БД
// Данные юзера передаются указателем на структуру
// Функция возвращает присвоенный в БД ID юзера или ошибку записи
func CreateUser(ctx context.Context, user *User) (userID int64, err error) {
	// Собрать запрос на вставку записи в таблицу
	builderQuery := sq.Insert(table).
		PlaceholderFormat(sq.Dollar).
		Columns(fieldName, fieldPassword, fieldEmail, fieldRole).
		Values(user.Name, user.Password, user.Email, user.Role).
		Suffix("RETURNING " + fieldID)

	query, args, err := builderQuery.ToSql()
	if err != nil {
		return 0, err
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// GetUser возвращает сведения о юзере по ID
func GetUser(ctx context.Context, userID int64) (user User, err error) {
	// Если юзера с таким ID не существует в базе, вернуть ошибку
	err = findUserByID(ctx, userID)
	if err != nil {
		return User{}, err
	}

	builderQuery := sq.Select(fieldName, fieldEmail, fieldRole, fieldCreatedAt, fieldUpdatedAt).
		From(table).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{fieldID: userID})

	query, args, err := builderQuery.ToSql()
	if err != nil {
		return User{}, err
	}

	var updatedAt sql.NullTime
	err = pool.QueryRow(ctx, query, args...).Scan(&user.Name, &user.Email, &user.Role, &user.CreatedAt, &updatedAt)
	if err != nil {
		return User{}, err
	}
	user.UpdatedAt = updatedAt.Time
	user.ID = userID

	return user, nil
}

// UpdateUser обновляет данные юзера в БД
func UpdateUser(ctx context.Context, userID int64, name stringValue, email stringValue, role desc.Role) (err error) {
	// Если юзера с таким ID не существует в базе, вернуть ошибку
	err = findUserByID(ctx, userID)
	if err != nil {
		return err
	}

	builderQuery := sq.Update(table).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{fieldID: userID})

	isValid := false
	if name != nil {
		builderQuery = builderQuery.Set(fieldName, name.Value)
		isValid = true
	}
	if email != nil {
		builderQuery = builderQuery.Set(fieldEmail, email.Value)
		isValid = true
	}
	if role != desc.Role_UNKNOWN {
		builderQuery = builderQuery.Set(fieldRole, role.String())
		isValid = true
	}
	if !isValid {
		return errors.New("there are no fields to update")
	}
	builderQuery = builderQuery.Set(fieldUpdatedAt, time.Now())

	query, args, err := builderQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser удаляет юзера в БД
func DeleteUser(ctx context.Context, userID int64) (err error) {
	// Если юзера с таким ID не существует в базе, вернуть ошибку
	err = findUserByID(ctx, userID)
	if err != nil {
		return err
	}

	builderQuery := sq.Delete(table).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{fieldID: userID})

	query, args, err := builderQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// findUserByID возвращает err, если юзера с указанным ID не существует в базе и nil в противном случае
func findUserByID(ctx context.Context, userID int64) error {
	// Создаем запрос с помощью Squirrel
	query, args, err := sq.Select(fieldID).
		PlaceholderFormat(sq.Dollar).
		From(table).
		Where(sq.Eq{fieldID: userID}).
		ToSql()
	if err != nil {
		return err
	}

	var ret int64
	err = pool.QueryRow(ctx, query, args...).Scan(&ret)
	if err != nil {
		return err
	}

	return nil
}
