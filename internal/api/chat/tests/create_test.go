package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	chatApi "github.com/algol-84/chat-server/internal/api/chat"
	"github.com/algol-84/chat-server/internal/model"
	"github.com/algol-84/chat-server/internal/service"
	serviceMocks "github.com/algol-84/chat-server/internal/service/mocks"
	desc "github.com/algol-84/chat-server/pkg/chat_v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	// mc - служебный объект minimock
	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	names := make([]string, 3)
	for i := 0; i < 3; i++ {
		names[i] = gofakeit.Name()
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id         = gofakeit.Int64()
		serviceErr = status.Errorf(codes.Internal, "user creation in DB returned with error")

		req = &desc.CreateRequest{
			Chat: &desc.Chat{
				Usernames: names,
			},
		}

		chat = &model.Chat{
			Usernames: names,
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
				mock.CreateMock.Expect(ctx, chat).Return(id, nil)
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
				mock.CreateMock.Expect(ctx, chat).Return(0, serviceErr)
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

			newID, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
