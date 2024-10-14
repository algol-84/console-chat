package main

import (
	"context"
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/algol-84/auth/internal/config"

	desc "github.com/algol-84/auth/pkg/user_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

const grpcPort = 50051

// const (
// 	dbDSN = "host=localhost port=54322 dbname=auth user=auth-user password=auth-password sslmode=disable"
// )

type server struct {
	desc.UnimplementedUserV1Server
	// dbWorker *DbWorker
}

// Create User
func (s *server) Create(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Request to create user received: %+v", req.Info)

	// TODO Check user exists in the database

	// if req.Info.Password != req.Info.PasswordConfirm {
	// 	// TODO return response with error
	// }

	// s.dbWorker.user.name = req.Info.Name
	// s.dbWorker.user.password = req.Info.Password
	// s.dbWorker.user.email = req.Info.Email
	// s.dbWorker.user.role = req.Info.Role.String()
	// s.dbWorker.user.created_at = time.Now()

	// userID, err := s.dbWorker.createUser()

	// if err != nil {
	// 	log.Println("DB error:", err)
	// }

	return &desc.CreateResponse{
		//	Id: userID,
	}, nil
}

// Get User info by ID
func (s *server) Get(_ context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("Request to get user %d", req.Id)

	// s.dbWorker.getUser(req.Id)
	// // Convert from string value to enum Role
	// var role desc.Role
	// if roleValue, exists := desc.Role_value[s.dbWorker.user.role]; exists {
	// 	role = desc.Role(roleValue)
	// } else {
	// 	role = desc.Role_UNKNOWN
	// }

	return &desc.GetResponse{
		User: &desc.GetUser{
			Id: req.Id,
			// Name:      s.dbWorker.user.name,
			// Email:     s.dbWorker.user.email,
			// Role:      role,
			// CreatedAt: timestamppb.New(s.dbWorker.user.created_at),
			// UpdatedAt: timestamppb.New(s.dbWorker.user.updated_at),
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

	// s.dbWorker.user.name = req.Info.Name.Value
	// s.dbWorker.user.email = req.Info.Email.Value
	// s.dbWorker.user.role = req.Info.Role.String()
	// s.dbWorker.updateUser(req.Id)

	return &emptypb.Empty{}, nil
}

// Delete User by ID
func (s *server) Delete(_ context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Request to delete user %d received", req.Id)

	// TODO Error handling if user not exists in the database

	//	s.dbWorker.deleteUser(req.Id)

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

	// dbWorker, err := NewDbWorker(dbDSN)
	// if err != nil {
	// 	log.Fatalf("Failed to create DbWorker: %v", err)
	// }
	// defer dbWorker.pool.Close()

	lis, err := net.Listen("tcp", grpcConfig.Address())

	//	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	userServer := &server{} //{dbWorker: dbWorker}
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, userServer)

	log.Printf("auth server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
