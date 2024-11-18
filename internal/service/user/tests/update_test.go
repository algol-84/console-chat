package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/algol-84/auth/internal/model"
	"github.com/algol-84/auth/internal/repository"
	repoMocks "github.com/algol-84/auth/internal/repository/mocks"
	auth "github.com/algol-84/auth/internal/service/user"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type authRepositoryMockFunc func(mc *minimock.Controller) repository.AuthRepository
	type cacheRepositoryMockFunc func(mc *minimock.Controller) repository.CacheRepository

	type args struct {
		ctx context.Context
		req *model.UserUpdate
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    int64
		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = "USER"

		repoErr  = fmt.Errorf("repo error")
		emptyErr = fmt.Errorf("update request is empty")

		req = &model.UserUpdate{
			ID: id,
			Name: model.StringValue{
				Value: name,
				Valid: true,
			},
			Email: model.StringValue{
				Value: email,
				Valid: true,
			},
			Role: model.StringValue{
				Value: role,
				Valid: true,
			},
		}

		reqFail = &model.UserUpdate{
			ID: id,
			Name: model.StringValue{
				Value: name,
				Valid: false,
			},
			Email: model.StringValue{
				Value: email,
				Valid: false,
			},
			Role: model.StringValue{
				Value: role,
				Valid: false,
			},
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name                string
		args                args
		want                int64
		err                 error
		authRepositoryMock  authRepositoryMockFunc
		cacheRepositoryMock cacheRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, req).Return(nil)
				return mock
			},
			cacheRepositoryMock: func(mc *minimock.Controller) repository.CacheRepository {
				mock := repoMocks.NewCacheRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, req.ID).Return(nil)
				return mock
			},
		},
		{
			name: "service PG error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: repoErr,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, req).Return(repoErr)
				return mock
			},
			cacheRepositoryMock: func(_ *minimock.Controller) repository.CacheRepository {
				// Если произошла ошибка при апдейте базы, то удалять юзера из кеша не нужно
				// функция не вызывается
				return nil
			},
		},
		{
			name: "service Redis error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: repoErr,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, req).Return(nil)
				return mock
			},
			cacheRepositoryMock: func(mc *minimock.Controller) repository.CacheRepository {
				mock := repoMocks.NewCacheRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, req.ID).Return(repoErr)
				return mock
			},
		},
		{
			name: "service empty request error case",
			args: args{
				ctx: ctx,
				req: reqFail,
			},
			err: emptyErr,
			authRepositoryMock: func(_ *minimock.Controller) repository.AuthRepository {
				return nil
			},
			cacheRepositoryMock: func(_ *minimock.Controller) repository.CacheRepository {
				return nil
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

			err := service.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
