package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/algol-84/auth/internal/model"
	"github.com/algol-84/auth/internal/repository"
	repoMocks "github.com/algol-84/auth/internal/repository/mocks"
	"github.com/algol-84/auth/internal/service/auth"
)

func TestGet(t *testing.T) {
	t.Parallel()
	// mc - служебный объект minimock
	type authRepositoryMockFunc func(mc *minimock.Controller) repository.AuthRepository
	type cacheRepositoryMockFunc func(mc *minimock.Controller) repository.CacheRepository

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Name()
		email     = gofakeit.Email()
		role      = "USER"
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		//	repoErr = fmt.Errorf("repo error")

		res = &model.User{
			ID:        id,
			Name:      name,
			Email:     email,
			Role:      role,
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{
				Time:  updatedAt,
				Valid: true,
			},
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name                string
		args                args
		want                *model.User
		err                 error
		authRepositoryMock  authRepositoryMockFunc
		cacheRepositoryMock cacheRepositoryMockFunc
	}{
		{
			name: "success case from cache",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: res,
			err:  nil,
			// При удачном запросе из кеша в базу не ходим
			authRepositoryMock: func(_ *minimock.Controller) repository.AuthRepository {
				return nil
			},
			cacheRepositoryMock: func(mc *minimock.Controller) repository.CacheRepository {
				mock := repoMocks.NewCacheRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(res, nil)
				return mock
			},
		},
		{
			name: "success case from PG",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: res,
			err:  nil,
			// В кэше юзер не найден, идем в базу
			cacheRepositoryMock: func(mc *minimock.Controller) repository.CacheRepository {
				mock := repoMocks.NewCacheRepositoryMock(mc)
				// Ожидаем получение юзера из кэша
				mock.GetMock.Expect(ctx, id).Return(nil, model.ErrorUserNotFound)
				// Ожидаем добавление юзера в кэш
				mock.CreateMock.Expect(ctx, res).Return(id, nil)
				return mock
			},
			// А в базе нашелся
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				// Ожидаем получение юзера из базы
				mock.GetMock.Expect(ctx, id).Return(res, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: nil,
			err:  model.ErrorUserNotFound,
			// Ожидаем ошибку что юзер не найден в базе
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, model.ErrorUserNotFound)
				return mock
			},
			// В кэше юзера тоже нет
			cacheRepositoryMock: func(mc *minimock.Controller) repository.CacheRepository {
				mock := repoMocks.NewCacheRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, model.ErrorUserNotFound)
				return mock
			},
		},
		{
			name: "service error case: cache internal error",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: nil,
			err:  model.ErrorCacheInternal,

			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				// Ожидаем получению юзера из базы
				mock.GetMock.Expect(ctx, id).Return(res, nil)
				return mock
			},

			cacheRepositoryMock: func(mc *minimock.Controller) repository.CacheRepository {
				mock := repoMocks.NewCacheRepositoryMock(mc)
				// Ожидаем, что юзер не найден в кэше
				mock.GetMock.Expect(ctx, id).Return(nil, model.ErrorUserNotFound)
				// Ожидаем добавление юзера в кэш, с возвратом ошибки записи
				mock.CreateMock.Expect(ctx, res).Return(id, model.ErrorCacheInternal)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authRepoMock := tt.authRepositoryMock(mc)
			cacheRepoMock := tt.cacheRepositoryMock(mc)
			service := auth.NewMockService(authRepoMock, cacheRepoMock)

			newID, err := service.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
