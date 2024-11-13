package auth

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/algol-84/auth/internal/repository"
	"github.com/algol-84/auth/internal/repository/auth/pg/converter"
	db "github.com/algol-84/platform_common/pkg/db"

	model "github.com/algol-84/auth/internal/model"
	modelRepo "github.com/algol-84/auth/internal/repository/auth/pg/model"
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
	db db.Client
}

// NewRepository конструктор
func NewRepository(db db.Client) repository.AuthRepository {
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

	q := db.Query{
		Name:     "auth_repository.Create",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID)
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

	q := db.Query{
		Name:     "auth_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
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

	q := db.Query{
		Name:     "auth_repository.Update",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
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

	q := db.Query{
		Name:     "auth_repository.Update",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID=%d not found", id)
	}

	return nil
}

func (r *repo) Find(ctx context.Context, username string) (*model.User, error) {
	builderQuery := sq.Select(fieldID, fieldName, fieldEmail, fieldRole, fieldPassword, fieldCreatedAt, fieldUpdatedAt).
		From(table).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{fieldName: username})

	query, args, err := builderQuery.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "auth_repository.Find",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}
