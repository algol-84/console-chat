package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

const table = "chat_user"

type DbWorker struct {
	pool *pgxpool.Pool
	ctx  context.Context
	user User
}

type User struct {
	name       string
	password   string
	email      string
	role       string
	created_at time.Time
	updated_at time.Time
}

func NewDbWorker(connString string) (*DbWorker, error) {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}
	return &DbWorker{pool: pool, ctx: ctx}, nil
}

// Function creates the user in the database
func (db *DbWorker) createUser() (userID int64, err error) {
	// Делаем запрос на вставку записи в таблицу note
	builderQuery := sq.Insert(table).
		PlaceholderFormat(sq.Dollar).
		Columns("name", "password", "email", "role", "created_at").
		Values(db.user.name, db.user.password, db.user.email, db.user.role, db.user.created_at).
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

// Function selects one user by ID from database
func (db *DbWorker) getUser(userID int64) (err error) {
	builderQuery := sq.Select("name", "email", "role", "created_at", "updated_at").
		From(table).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": userID})

	query, args, err := builderQuery.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var updated_at sql.NullTime
	err = db.pool.QueryRow(db.ctx, query, args...).Scan(&db.user.name, &db.user.email, &db.user.role, &db.user.created_at, &updated_at)
	if err != nil {
		log.Fatalf("failed to select notes: %v", err)
	}

	log.Printf("Select user info from DB: %v", db.user)
	return nil
}

// Function updates the user by ID in the database
func (db *DbWorker) updateUser(userID int64) (err error) {
	builderQuery := sq.Update(table).
		PlaceholderFormat(sq.Dollar).
		Set("name", db.user.name).
		Set("email", db.user.email).
		Set("role", db.user.role).
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

// Function removes the user by ID from the database
func (db *DbWorker) deleteUser(userID int64) (err error) {
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
