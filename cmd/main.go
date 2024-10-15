package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/algol-84/auth/internal/config"
	pg "github.com/algol-84/auth/internal/pg_auth"

	desc "github.com/algol-84/auth/pkg/user_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedUserV1Server
	dbWorker *pg.DbWorker
}

// Create User
func (s *server) Create(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Request to create user received: %+v", req.Info)

	// TODO Check user exists in the database

	// if req.Info.Password != req.Info.PasswordConfirm {
	// 	// TODO return response with error
	// }

	s.dbWorker.User.Name = req.Info.Name
	s.dbWorker.User.Password = req.Info.Password
	s.dbWorker.User.Email = req.Info.Email
	s.dbWorker.User.Role = req.Info.Role.String()
	s.dbWorker.User.CreatedAt = time.Now()

	userID, err := s.dbWorker.CreateUser()

	if err != nil {
		log.Println("DB error:", err)
	}

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}

// Get User info by ID
func (s *server) Get(_ context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("Request to get user %d", req.Id)

	err := s.dbWorker.GetUser(req.Id)
	if err != nil {
		log.Println(err)
	}
	// Convert from string value to enum Role
	var role desc.Role
	if roleValue, exists := desc.Role_value[s.dbWorker.User.Role]; exists {
		role = desc.Role(roleValue)
	} else {
		role = desc.Role_UNKNOWN
	}

	return &desc.GetResponse{
		User: &desc.GetUser{
			Id:        req.Id,
			Name:      s.dbWorker.User.Name,
			Email:     s.dbWorker.User.Email,
			Role:      role,
			CreatedAt: timestamppb.New(s.dbWorker.User.CreatedAt),
			UpdatedAt: timestamppb.New(s.dbWorker.User.UpdatedAt),
		},
	}, nil
}

// Update User info
func (s *server) Update(_ context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("Request to update user %d received", req.Id)

	if req.Info.Name == nil || req.Info.Email == nil {
		// TODO return error response
		return &emptypb.Empty{}, nil
	}

	s.dbWorker.User.Name = req.Info.Name.Value
	s.dbWorker.User.Email = req.Info.Email.Value
	s.dbWorker.User.Role = req.Info.Role.String()
	err := s.dbWorker.UpdateUser(req.Id)
	if err != nil {
		log.Println(err)
	}

	return &emptypb.Empty{}, nil
}

// Delete User by ID
func (s *server) Delete(_ context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Request to delete user %d received", req.Id)

	// TODO Error handling if user not exists in the database

	err := s.dbWorker.DeleteUser(req.Id)
	if err != nil {
		log.Println(err)
	}

	return &emptypb.Empty{}, nil
}

func main() {
	log.Printf("Run auth-server...")

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

	log.Println("-->>", pgConfig)

	ctx := context.Background()
	dbWorker, err := pg.NewDbWorker(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("Failed to create DbWorker: %v", err)
	}
	defer dbWorker.Close()

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	userServer := &server{dbWorker: dbWorker}
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, userServer)

	log.Printf("auth server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
