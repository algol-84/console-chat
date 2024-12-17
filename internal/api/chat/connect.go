package chat

import (
	desc "github.com/algol-84/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ConnectChat осуществляет подключение к существующему чату
func (i *Implementation) ConnectChat(req *desc.ConnectChatRequest, stream desc.ChatV1_ConnectChatServer) error {
	// Проверяем, что чат существует
	i.mxChannel.RLock()
	chatChan, ok := i.channels[req.GetChatId()]
	i.mxChannel.RUnlock()

	if !ok {
		return status.Errorf(codes.NotFound, "chat not found")
	}

	i.mxChat.Lock()

	if _, okChat := i.chats[req.GetChatId()]; !okChat {
		i.chats[req.GetChatId()] = &Chat{
			streams: make(map[string]desc.ChatV1_ConnectChatServer),
		}
	}
	i.mxChat.Unlock()

	i.chats[req.GetChatId()].m.Lock()
	i.chats[req.GetChatId()].streams[req.GetUsername()] = stream
	i.chats[req.GetChatId()].m.Unlock()

	for {
		select {
		// Ожидаем данных в канале
		case msg, okCh := <-chatChan:
			if !okCh {
				return nil
			}
			// Разослать сообщение всем участникам чата
			for _, st := range i.chats[req.GetChatId()].streams {
				// Отправить сообщение в стрим
				if err := st.Send(msg); err != nil {
					return err
				}
			}

			// Ожидаем завершения контекста
		case <-stream.Context().Done():
			i.chats[req.GetChatId()].m.Lock()
			// Стрим закрыт, удаляем пользователя из чата
			delete(i.chats[req.GetChatId()].streams, req.GetUsername())
			i.chats[req.GetChatId()].m.Unlock()
			return nil
		}
	}
}
