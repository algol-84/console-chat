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
}

// Create User
func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	// Проверить корректность всех полей запроса
	if req.Password != req.PasswordConfirm {
		return nil, status.Errorf(codes.InvalidArgument, "fields \"password\" and \"password_confirm\" don't match")
	}

	if req.Role != desc.Role_ADMIN && req.Role != desc.Role_USER {
		return nil, status.Errorf(codes.InvalidArgument, "user role must be ADMIN or USER")
	}

	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user name must be not empty")
	}

	var user pg.User
	user.Name = req.Name
	user.Email = req.Email
	user.Password = req.Password
	user.Role = req.Role.String()

	userID, err := pg.CreateUser(ctx, &user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "user creation in DB returned with error: %s", err)
	}

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}

func convertStringToRole(roleStr string) desc.Role {
	if roleValue, exists := desc.Role_value[roleStr]; exists {
		return desc.Role(roleValue)
	}
	return desc.Role_UNKNOWN
}

// Get User info by ID
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := pg.GetUser(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "the request for user data in the DB returned with error: %s", err)
	}

	return &desc.GetResponse{
		Id:        req.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      convertStringToRole(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

// Update User info
func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	if req.Name == nil && req.Email == nil && req.Role == desc.Role_UNKNOWN {
		return nil, status.Errorf(codes.InvalidArgument, "there are no fields to update")
	}

	err := pg.UpdateUser(ctx, req.Id, req.Name, req.Email, req.Role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "updating user in the DB returned with an error: %s", err)
	}

	return &emptypb.Empty{}, nil
}

// Delete User by ID
func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := pg.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "removing user from the DB returned with an error: %s", err)
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

	ctx := context.Background()
	err = pg.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("Failed to create DbWorker: %v", err)
	}
	defer pg.Close()

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	userServer := &server{}
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, userServer)

	log.Printf("auth server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
