package main

import (
	"context"
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/algol-84/auth/internal/config"
	"github.com/algol-84/auth/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/algol-84/auth/internal/repository/auth"

	desc "github.com/algol-84/auth/pkg/user_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedUserV1Server
	authRepository repository.AuthRepository
}

// Create обрабатывает GRPC запросы на создание нового юзера
func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	userID, err := s.authRepository.Create(ctx, req.User)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "user creation in DB returned with error: %s", err)
	}

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}

// Get обрабатывает GRPC запросы на получение данных пользователя
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	userInfo, err := s.authRepository.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "the request for user data in the DB returned with error: %s", err)
	}

	return &desc.GetResponse{
		UserInfo: userInfo,
	}, nil
}

// Update обрабатывает GRPC запросы на обновление данных пользователя
func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	err := s.authRepository.Update(ctx, req.UserUpdate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "updating user in the DB returned with an error: %s", err)
	}

	return &emptypb.Empty{}, nil
}

// Delete обрабатывает GRPC запросы на удаление пользователя
func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := s.authRepository.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "removing user from the DB returned with an error: %s", err)
	}

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

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	authRepo := auth.NewRepository(pool)

	userServer := &server{authRepository: authRepo}
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, userServer)

	log.Printf("auth server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
