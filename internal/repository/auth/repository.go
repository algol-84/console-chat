package auth

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/algol-84/auth/internal/repository"
	"github.com/algol-84/auth/internal/repository/auth/converter"

	model "github.com/algol-84/auth/internal/model"
	modelRepo "github.com/algol-84/auth/internal/repository/auth/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

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

type repo struct {
	db *pgxpool.Pool
}

// NewRepository конструктор
func NewRepository(db *pgxpool.Pool) repository.AuthRepository {
	return &repo{db: db}
}

// Create создает нового юзера в БД
// Данные юзера передаются указателем на структуру protobuf
// Функция возвращает присвоенный в БД ID юзера или ошибку записи
func (r *repo) Create(ctx context.Context, user *model.User) (int64, error) {
	var userID int64
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

	err = r.db.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// Get возвращает информацию о юзере по ID
func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	builderQuery := sq.Select(fieldID, fieldName, fieldEmail, fieldRole, fieldCreatedAt, fieldUpdatedAt).
		From(table).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{fieldID: id})

	query, args, err := builderQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var user modelRepo.User
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

// Update обновляет данные юзера в БД
func (r *repo) Update(ctx context.Context, user *model.UserUpdate) error {
	builderQuery := sq.Update(table).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{fieldID: user.ID})

	if user.Name.Valid {
		builderQuery = builderQuery.Set(fieldName, user.Name.Value)
	}
	if user.Email.Valid {
		builderQuery = builderQuery.Set(fieldEmail, user.Email.Value)
	}
	if user.Role.Valid {
		builderQuery = builderQuery.Set(fieldRole, user.Role.Value)
	}
	builderQuery = builderQuery.Set(fieldUpdatedAt, time.Now())

	query, args, err := builderQuery.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID=%d not found", user.ID)
	}

	return nil
}

// Delete удаляет юзера из БД
func (r *repo) Delete(ctx context.Context, id int64) error {
	builderQuery := sq.Delete(table).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{fieldID: id})

	query, args, err := builderQuery.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID=%d not found", id)
	}

	return nil
}
