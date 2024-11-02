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

	chatApi "github.com/algol-84/chat-server/internal/api/chat"
	"github.com/algol-84/chat-server/internal/service"
	serviceMocks "github.com/algol-84/chat-server/internal/service/mocks"
	desc "github.com/algol-84/chat-server/pkg/chat_v1"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	// mc - служебный объект minimock
	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *desc.DeleteRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id         = gofakeit.Int64()
		serviceErr = status.Errorf(codes.Internal, "removing user from the DB returned with an error")

		req = &desc.DeleteRequest{
			Id: id,
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
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
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
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatServiceMock := tt.chatServiceMock(mc)
			api := chatApi.NewImplementation(chatServiceMock)

			newID, err := api.Delete(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
