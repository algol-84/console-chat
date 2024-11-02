package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/algol-84/auth/internal/api/auth"
	"github.com/algol-84/auth/internal/model"
	"github.com/algol-84/auth/internal/service"
	serviceMocks "github.com/algol-84/auth/internal/service/mocks"
	desc "github.com/algol-84/auth/pkg/user_v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	// mc - служебный объект minimock
	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       int64
		name     = gofakeit.Name()
		password = gofakeit.Password(true, false, false, false, false, 10)
		email    = gofakeit.Email()
		role     = "USER"

		serviceErr = status.Errorf(codes.Internal, "user creation in DB returned with error")

		req = &desc.CreateRequest{
			User: &desc.User{
				Name:            name,
				Password:        password,
				PasswordConfirm: password,
				Email:           email,
				Role:            desc.Role_USER,
			},
		}

		user = &model.User{
			Name:            name,
			Password:        password,
			PasswordConfirm: password,
			Email:           email,
			Role:            role,
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)
	defer t.Cleanup(mc.Finish)

	// Тесты объявлены слайсом структур
	tests := []struct {
		name            string               // Имя теста
		args            args                 // Аргументы
		want            *desc.CreateResponse // тип ожидаемого результата
		err             error                // ожидаемая ошибка
		authServiceMock authServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.CreateMock.Expect(ctx, user).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.CreateMock.Expect(ctx, user).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authServiceMock := tt.authServiceMock(mc)
			api := auth.NewImplementation(authServiceMock)

			newID, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
