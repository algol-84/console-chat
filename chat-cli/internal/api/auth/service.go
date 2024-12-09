package auth

import (
	"fmt"
	"log"

	descAuth "github.com/algol-84/chat-cli/pkg/auth_v1"
	user "github.com/algol-84/chat-cli/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthImpl struct {
	userClient user.UserV1Client
	authClient descAuth.AuthV1Client
}

func NewAuthClient() *AuthImpl {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", authHost, authPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial GRPC client: %v", err)
	}

	userClient := user.NewUserV1Client(conn)
	authClient := descAuth.NewAuthV1Client(conn)

	return &AuthImpl{
		userClient: userClient,
		authClient: authClient,
	}
}
