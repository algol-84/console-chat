package root

import (
	"flag"
	"log"

	"github.com/algol-84/chat-cli/internal/api/auth"
	"github.com/algol-84/chat-cli/internal/api/chat"
	"github.com/algol-84/chat-cli/internal/config"
	"github.com/robfig/cron"
)

var configPath string

// ServiceProvider хранит все объекты приложения, как интерфейсы или ссылки на структуры
type ClientProvider struct {
	grpcConfig   config.GRPCConfig
	authClient   *auth.AuthImpl
	chatClient   *chat.ChatImpl
	refreshToken string
	accessToken  string
}

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func NewClientProvider() (*ClientProvider, error) {

	flag.Parse()
	err := config.Load(configPath)
	if err != nil {
		return nil, err
	}

	cfg, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %s", err.Error())
	}

	authClient, err := auth.NewAuthClient(cfg.AuthServiceAddress())
	if err != nil {
		log.Fatalf("failed to connect to auth service: %s", err.Error())
	}

	chatClient, err := chat.NewChatClient(cfg.ChatServiceAddress())
	if err != nil {
		log.Fatalf("failed to connect to chat service: %s", err.Error())
	}

	return &ClientProvider{
		grpcConfig: cfg,
		authClient: authClient,
		chatClient: chatClient,
	}, nil
}

func (cl *ClientProvider) Login(username string, password string) {
	// Run Cron to periodic refresh tokens
	c := cron.New()
	err := c.AddFunc("@every 15m",
		func() {
			cl.updateAccessToken(cl.refreshToken)
			log.Println("new access token:", cl.accessToken)
		})
	if err != nil {
		log.Fatalf("internal cron error")
	}
	c.Start()

	// Логинимся и получаем refresh token
	token, err := cl.authClient.Login(username, password)
	if err != nil {
		log.Fatalf("fail to login user: %s", err.Error())
	}
	cl.refreshToken = token
	// Выписываем access token
	log.Println("refresh>> ", cl.refreshToken)
	token, err = cl.authClient.GetAccessToken(cl.refreshToken)
	if err != nil {
		log.Fatalf("fail to get access token: %s", err.Error())
	}
	cl.accessToken = token
}

func (cl *ClientProvider) updateAccessToken(refreshToken string) {
	token, err := cl.authClient.GetAccessToken(refreshToken)
	if err != nil {
		log.Fatalf("fail to get access token: %s", err.Error())
	}
	log.Println("new access token:", token)

	cl.accessToken = token
}
