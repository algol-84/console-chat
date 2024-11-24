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
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
)

const (
	// authPrefix префикс добавляется к токену для идентификации используемого метода аутентификации.
	// В случае JWT принято добавлять Bearer
	authPrefix      = "Bearer "
	authServicePort = 50051
)

type authConn struct {
	client descAccess.AccessV1Client
}

// NewAuthConnection создает коннект к Auth сервису
func NewAuthConnection() (*authConn, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf(":%d", authServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		return nil, errors.New("failed to dial auth service")
	}

	client := descAccess.NewAccessV1Client(conn)

	return &authConn{client: client}, nil
}

// AuthInterceptor интерцептор дергает ручку Check сервиса Auth, рефреш токен передается в контексте
func (a *authConn) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, info.FullMethod)
	defer span.Finish()

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

	//В контексте лежит ключ авторизации
	if ok && len(authHeader) > 0 {
		authEndpoint, ok := md["endpoint"]
		if !ok || len(authEndpoint) == 0 {
			return nil, errors.New("endpoint is not provided")
		}

		ctxOut := context.Background()
		md = metadata.New(map[string]string{"Authorization": authHeader[0]})
		ctxOut = metadata.NewOutgoingContext(ctxOut, md)

		_, err := a.client.Check(ctxOut, &descAccess.CheckRequest{
			EndpointAddress: authEndpoint[0],
		})
		if err != nil {
			return nil, fmt.Errorf("failed to send check request to auth service: %v", err)
		}
	}

	return handler(ctx, req)
}
