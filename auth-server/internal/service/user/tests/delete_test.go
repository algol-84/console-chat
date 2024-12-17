package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/algol-84/auth/internal/logger"
	"github.com/algol-84/auth/internal/repository"
	repoMocks "github.com/algol-84/auth/internal/repository/mocks"
	auth "github.com/algol-84/auth/internal/service/user"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	type authRepositoryMockFunc func(mc *minimock.Controller) repository.AuthRepository
	type cacheRepositoryMockFunc func(mc *minimock.Controller) repository.CacheRepository

	logger.MockInit()

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id      = gofakeit.Int64()
		repoErr = fmt.Errorf("repo error")
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
				req: id,
			},
			err: nil,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
			cacheRepositoryMock: func(mc *minimock.Controller) repository.CacheRepository {
				mock := repoMocks.NewCacheRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: id,
			},
			err: repoErr,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(repoErr)
				return mock
			},
			cacheRepositoryMock: func(mc *minimock.Controller) repository.CacheRepository {
				mock := repoMocks.NewCacheRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(repoErr)
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

			err := service.Delete(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
