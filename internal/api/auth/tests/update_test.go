package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/algol-84/auth/internal/api/auth"
	"github.com/algol-84/auth/internal/model"
	"github.com/algol-84/auth/internal/service"
	serviceMocks "github.com/algol-84/auth/internal/service/mocks"
	desc "github.com/algol-84/auth/pkg/user_v1"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	// mc - служебный объект minimock
	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type args struct {
		ctx context.Context
		req *desc.UpdateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = "USER"

		serviceErr = status.Errorf(codes.Internal, "updating user in the DB returned with an error")

		req = &desc.UpdateRequest{
			UserUpdate: &desc.UserUpdate{
				Id: id,
				Name: &wrapperspb.StringValue{
					Value: name,
				},
				Email: &wrapperspb.StringValue{
					Value: email,
				},
				Role: desc.Role_USER,
			},
		}

		user = &model.UserUpdate{
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

		res = &empty.Empty{}
	)
	defer t.Cleanup(mc.Finish)

	// Тесты объявлены слайсом структур
	tests := []struct {
		name            string         // Имя теста
		args            args           // Аргументы
		want            *emptypb.Empty // тип ожидаемого результата
		err             error          // ожидаемая ошибка
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
				mock.UpdateMock.Expect(ctx, user).Return(nil)
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
				mock.UpdateMock.Expect(ctx, user).Return(serviceErr)
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

			newID, err := api.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
