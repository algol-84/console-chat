package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/algol-84/auth/internal/api/auth"
	"github.com/algol-84/auth/internal/converter"
	"github.com/algol-84/auth/internal/model"
	"github.com/algol-84/auth/internal/service"
	serviceMocks "github.com/algol-84/auth/internal/service/mocks"
	desc "github.com/algol-84/auth/pkg/user_v1"
)

func TestGet(t *testing.T) {
	t.Parallel()
	// mc - служебный объект minimock
	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Name()
		email     = gofakeit.Email()
		role      = desc.Role_USER.String()
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		serviceErr = status.Errorf(codes.Internal, "the request for user data in the DB returned with error")

		req = &desc.GetRequest{
			Id: id,
		}

		serviceRes = &model.User{
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

		res = &desc.GetResponse{
			UserInfo: &desc.UserInfo{
				Id:        id,
				Name:      name,
				Email:     email,
				Role:      converter.StringToRole(role),
				CreatedAt: timestamppb.New(createdAt),
				UpdatedAt: timestamppb.New(updatedAt),
			},
		}
	)
	defer t.Cleanup(mc.Finish)

	// Тесты объявлены слайсом структур
	tests := []struct {
		name            string            // Имя теста
		args            args              // Аргументы
		want            *desc.GetResponse // тип ожидаемого результата
		err             error             // ожидаемая ошибка
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
				mock.GetMock.Expect(ctx, id).Return(serviceRes, nil)
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
				mock.GetMock.Expect(ctx, id).Return(&model.User{}, serviceErr)
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

			newID, err := api.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
