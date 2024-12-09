package chat

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/algol-84/chat-cli/pkg/chat_v1"
	"github.com/fatih/color"
	//"github.com/brianvoe/gofakeit"
)

const (
	chatHost = "127.0.0.1"
	chatPort = 50052
)

type ChatImpl struct {
	chatClient desc.ChatV1Client
}

func NewChatService() *ChatImpl {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", chatHost, chatPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial GRPC client: %v", err)
	}

	chatClient := desc.NewChatV1Client(conn)

	return &ChatImpl{chatClient: chatClient}
}

// Create new chat and return UUID
func (i *ChatImpl) Create(ctx context.Context) (string, error) {
	res, err := i.chatClient.CreateChat(ctx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}

	return res.Id, nil
}

func (i *ChatImpl) Connect(ctx context.Context, chatID string, username string) error {
	stream, err := i.chatClient.ConnectChat(ctx, &desc.ConnectChatRequest{
		ChatId:   chatID,
		Username: username,
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			message, errRecv := stream.Recv()
			if errRecv == io.EOF {
				return
			}
			if errRecv != nil {
				log.Println("failed to receive message from stream: ", errRecv)
				return
			}

			log.Printf("[%v] - [from: %s]: %s\n",
				color.YellowString(message.GetCreatedAt().AsTime().Format(time.RFC3339)),
				color.BlueString(message.GetFrom()),
				message.GetText(),
			)
		}
	}()

	for {
		scanner := bufio.NewScanner(os.Stdin)
		var lines strings.Builder

		for {
			scanner.Scan()
			line := scanner.Text()
			if len(line) == 0 {
				break
			}

			lines.WriteString(line)
			lines.WriteString("\n")
		}

		err = scanner.Err()
		if err != nil {
			log.Println("failed to scan message: ", err)
		}

		_, err = i.chatClient.SendMessage(ctx, &desc.SendMessageRequest{
			ChatId: chatID,
			Message: &desc.Message{
				From:      username,
				Text:      lines.String(),
				CreatedAt: timestamppb.Now(),
			},
		})
		if err != nil {
			log.Println("failed to send message: ", err)
			return err
		}
	}
}

func (i *ChatImpl) Delete(ctx context.Context, chatID string) error {
	_, err := i.chatClient.DeleteChat(ctx, &desc.DeleteChatRequest{
		Id: chatID,
	})
	if err != nil {
		return err
	}

	return nil
}
