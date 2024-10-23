package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/algol-84/chat-server/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pg "github.com/algol-84/chat-server/internal/pg"
	desc "github.com/algol-84/chat-server/pkg/chat_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedChatV1Server
}

// Create обработчик запроса на создание нового чата
func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	chatID, err := pg.CreateChat(ctx, req.Usernames)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "chat creation in DB returned with error: %s", err)
	}

	return &desc.CreateResponse{
		Id: chatID,
	}, nil
}

// Delete chat by ID
func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := pg.DeleteChat(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "chat deletion in DB returned with error: %s", err)
	}

	return &emptypb.Empty{}, nil
}

// Send message to chat
func (s *server) SendMessage(_ context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("Message %s received from %s at %v", req.Text, req.From, req.Timestamp.AsTime())

	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
	// Считываем переменные окружения
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v %v", err, configPath)
		return
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
		return
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	ctx := context.Background()
	err = pg.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to create DbWorker: %v", err)
	}
	defer pg.Close()

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, &server{})

	log.Printf("chat server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
