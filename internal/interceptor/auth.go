package interceptor

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	// TODO взять access_v1 из сервиса auth. Но пока access ручки находятся в пулл-реквесте и не смержены с основной веткой
	//descAccess "github.com/algol-84/auth/pkg/access_v1"
	descAccess "github.com/algol-84/chat-server/pkg/access_v1"
)

const (
	// authPrefix префикс добавляется к токену для идентификации используемого метода аутентификации.
	// В случае JWT принято добавлять Bearer
	authPrefix      = "Bearer "
	authServicePort = 50051
)

// AuthInterceptor интерцептор дергает ручку Check сервиса Auth, рефреш токен передается в контексте
func AuthInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("authorization header is not provided")
	}

	authEndpoint, ok := md["endpoint"]
	if !ok || len(authEndpoint) == 0 {
		return nil, errors.New("endpoint is not provided")
	}

	ctxOut := context.Background()
	md = metadata.New(map[string]string{"Authorization": authHeader[0]})
	ctxOut = metadata.NewOutgoingContext(ctxOut, md)

	conn, err := grpc.Dial(
		fmt.Sprintf(":%d", authServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, errors.New("failed to dial auth service")
	}

	cl := descAccess.NewAccessV1Client(conn)
	_, err = cl.Check(ctxOut, &descAccess.CheckRequest{
		EndpointAddress: authEndpoint[0],
	})
	if err != nil {
		return nil, errors.New("failed to send check request to auth service")
	}

	return handler(ctx, req)
}
