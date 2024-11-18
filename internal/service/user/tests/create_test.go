package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/algol-84/auth/internal/client/kafka"
	kafkaMocks "github.com/algol-84/auth/internal/client/kafka/mocks"
	"github.com/algol-84/auth/internal/model"
	"github.com/algol-84/auth/internal/repository"
	repoMocks "github.com/algol-84/auth/internal/repository/mocks"
	auth "github.com/algol-84/auth/internal/service/user"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type authRepositoryMockFunc func(mc *minimock.Controller) repository.AuthRepository
	type kafkaProducerMockFunc func(mc *minimock.Controller) kafka.Producer

	type args struct {
		ctx  context.Context
		req  *model.User
		data []byte
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       int64
		name     = gofakeit.Name()
		password = gofakeit.Password(true, false, false, false, false, 10)
		email    = gofakeit.Email()
		role     = "USER"

		repoErr = fmt.Errorf("repo error")

		req = &model.User{
			Name:            name,
			Password:        password,
			PasswordConfirm: password,
			Email:           email,
			Role:            role,
		}

		data, _ = json.Marshal(req)
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		authRepositoryMock authRepositoryMockFunc
		kafkaProducerMock  kafkaProducerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:  ctx,
				req:  req,
				data: data,
			},
			want: id,
			err:  nil,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
			kafkaProducerMock: func(mc *minimock.Controller) kafka.Producer {
				mock := kafkaMocks.NewProducerMock(mc)
				mock.ProduceMock.Expect(ctx, data).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx:  ctx,
				req:  req,
				data: data,
			},
			want: 0,
			err:  repoErr,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(0, repoErr)
				return mock
			},
			kafkaProducerMock: func(_ *minimock.Controller) kafka.Producer {
				// Если запись в базу вернулась с ошибкой, то запись в кафку не происходит
				return nil
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authRepoMock := tt.authRepositoryMock(mc)
			kafkaMock := tt.kafkaProducerMock(mc)
			service := auth.NewMockService(authRepoMock, kafkaMock)

			newID, err := service.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
