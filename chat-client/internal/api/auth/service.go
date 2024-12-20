package auth

import (
	"fmt"

	descAuth "github.com/algol-84/chat-cli/pkg/auth_v1"
	user "github.com/algol-84/chat-cli/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthImpl struct {
	userClient user.UserV1Client
	authClient descAuth.AuthV1Client
}

func NewAuthClient(address string) (*AuthImpl, error) {
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial GRPC client: %v", err)
	}

	userClient := user.NewUserV1Client(conn)
	authClient := descAuth.NewAuthV1Client(conn)

	return &AuthImpl{
		userClient: userClient,
		authClient: authClient,
	}, nil
}
