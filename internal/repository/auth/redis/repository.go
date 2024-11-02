package pg

import (
	"context"
	"errors"
	"strconv"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/algol-84/auth/internal/client/cache"
	"github.com/algol-84/auth/internal/model"
	"github.com/algol-84/auth/internal/repository"
	"github.com/algol-84/auth/internal/repository/auth/redis/converter"
	modelRepo "github.com/algol-84/auth/internal/repository/auth/redis/model"
)

type repo struct {
	cl cache.RedisClient
}

// NewRepository конструктор репо слоя редиса
func NewRepository(cl cache.RedisClient) repository.CacheRepository {
	return &repo{cl: cl}
}

// Create создает юзера в кэше
func (r *repo) Create(ctx context.Context, user *model.User) (int64, error) {
	usr := converter.ToRepoFromUser(user)
	idStr := strconv.FormatInt(user.ID, 10)
	err := r.cl.HashSet(ctx, idStr, usr)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

// Get запрашивает юзера из кэша
func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	idStr := strconv.FormatInt(id, 10)
	values, err := r.cl.HGetAll(ctx, idStr)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, model.ErrorUserNotFound
	}

	var user modelRepo.User
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

// func (r *repo) Update(_ context.Context, _ *model.UserUpdate) error {
// 	return nil
// }

func (r *repo) Delete(ctx context.Context, id int64) error {
	idStr := strconv.FormatInt(id, 10)

	err := r.cl.Del(ctx, idStr)
	if err != nil {
		//return err
		return errors.New("deletion error")
	}

	return nil
}
