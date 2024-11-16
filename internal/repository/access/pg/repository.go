package pg

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/algol-84/auth/internal/repository"
	db "github.com/algol-84/platform_common/pkg/db"
)

// Представление БД сервиса Access
const table = "chat_permissions"
const (
	fieldEndpoint = "endpoint"
	fieldRole     = "role"
)

type repo struct {
	db db.Client
}

// NewRepository конструктор
func NewRepository(db db.Client) repository.AccessRepository {
	return &repo{db: db}
}

type permissions struct {
	Endpoint string `db:"endpoint"`
	Role     string `db:"role"`
}

// Get возвращает права доступа для ручки endpoint
func (r *repo) Get(ctx context.Context) (map[string]string, error) {
	builderQuery := sq.Select(fieldEndpoint, fieldRole).
		From(table).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builderQuery.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "access_repository.Get",
		QueryRaw: query,
	}

	var p []permissions
	err = r.db.DB().ScanAllContext(ctx, &p, q, args...)

	if err != nil {
		return nil, err
	}

	roles := make(map[string]string)
	for _, value := range p {
		roles[value.Endpoint] = value.Role
	}

	return roles, nil
}
