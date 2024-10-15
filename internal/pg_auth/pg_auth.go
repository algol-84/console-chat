package pg_auth

import (
	"context"
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

const table = "chat_user"

// DbWorker содержит все переменные работы с БД auth
type DbWorker struct {
	pool *pgxpool.Pool
	ctx  context.Context
	User UserInfo
}

// UserInfo содержит все сведения о юзере из БД auth
type UserInfo struct {
	Name      string
	Password  string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewDbWorker создает пул подключений к БД
func NewDbWorker(ctx context.Context, connString string) (*DbWorker, error) {
	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}
	log.Println("Connection to the database was successful")
	return &DbWorker{pool: pool, ctx: ctx}, nil
}

// Close закрывает пул соединений с БД
func (db *DbWorker) Close() {
	db.pool.Close()
}

// CreateUser создает нового юзера в БД
func (db *DbWorker) CreateUser() (userID int64, err error) {
	// Делаем запрос на вставку записи в таблицу note
	builderQuery := sq.Insert(table).
		PlaceholderFormat(sq.Dollar).
		Columns("name", "password", "email", "role", "created_at").
		Values(db.User.Name, db.User.Password, db.User.Email, db.User.Role, db.User.CreatedAt).
		Suffix("RETURNING id")

	query, args, err := builderQuery.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
		return 0, err
	}

	err = db.pool.QueryRow(db.ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
		return 0, err
	}

	log.Printf("inserted user with id: %d", userID)
	return userID, nil
}

// GetUser возвращает сведения о юзере по ID
func (db *DbWorker) GetUser(userID int64) (err error) {
	builderQuery := sq.Select("name", "email", "role", "created_at", "updated_at").
		From(table).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": userID})

	query, args, err := builderQuery.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var updatedAt sql.NullTime
	err = db.pool.QueryRow(db.ctx, query, args...).Scan(&db.User.Name, &db.User.Email, &db.User.Role, &db.User.CreatedAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to select notes: %v", err)
	}

	// Перевод из sql.NullTime в time.Time
	if updatedAt.Valid {
		db.User.UpdatedAt = updatedAt.Time
	} else {
		db.User.UpdatedAt = time.Time{}
	}

	log.Printf("Select user info from DB: %v", db.User)
	return nil
}

// UpdateUser обновляет данные юзера в БД
func (db *DbWorker) UpdateUser(userID int64) (err error) {
	builderQuery := sq.Update(table).
		PlaceholderFormat(sq.Dollar).
		Set("name", db.User.Name).
		Set("email", db.User.Email).
		Set("role", db.User.Role).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": userID})

	query, args, err := builderQuery.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
		return err
	}

	res, err := db.pool.Exec(db.ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update user: %v", err)
		return err
	}

	log.Printf("DB updated %d rows", res.RowsAffected())
	return nil
}

// DeleteUser удаляет юзера в БД
func (db *DbWorker) DeleteUser(userID int64) (err error) {
	builderQuery := sq.Delete(table).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": userID})

	query, args, err := builderQuery.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
		return err
	}

	_, err = db.pool.Exec(db.ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to delete user: %v", err)
		return err
	}

	log.Printf("deleted user with id: %d", userID)
	return nil
}
