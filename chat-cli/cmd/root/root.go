package root

import (
	"context"
	"log"
	"os"

	"github.com/algol-84/chat-cli/internal/model"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

const (
	// authPrefix префикс добавляется к токену для идентификации используемого метода аутентификации.
	// В случае JWT принято добавлять Bearer
	authPrefix          = "Bearer "
	createChatEndpoint  = "/chat_v1.ChatV1/CreateChat"
	connectChatEndpoint = "/chat_v1.ChatV1/ConnectChat"
	sendMessageEndpoint = "/chat_v1.ChatV1/SendMessage"
)

var (
	userID          int64
	username        string
	email           string
	password        string
	passwordConfirm string
	role            string
	chatID          string
)

var cl *ClientProvider

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chat-client-app",
	Short: "chat client cli util",
}

var createCmd = &cobra.Command{
	Use:   "create",      // команда в терминале
	Short: "create user", // короткое описание команды для хел
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete user",
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update user",
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get user",
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect to chat",
}

// Create new chat
var createChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Create new chatroom",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		// Authorized user and get JWT tokens
		cl.Login(username, password)
		// Save token and endpoint to context
		// Note: only users with admin role can create the chat
		md := metadata.New(map[string]string{"Authorization": authPrefix + cl.accessToken, "Endpoint": createChatEndpoint})
		ctx = metadata.NewOutgoingContext(ctx, md)
		// Create chat
		chatId, err := cl.chatClient.Create(ctx)
		if err != nil {
			log.Fatalf("failed to create chat: %s\n", err.Error())
		}

		log.Printf("chat was created with id: %s\n", chatId)

		ctx = context.Background()
		md = metadata.New(map[string]string{"Authorization": authPrefix + cl.accessToken, "Endpoint": connectChatEndpoint})
		ctx = metadata.NewOutgoingContext(ctx, md)

		// Connect to created chat
		err = cl.chatClient.Connect(ctx, chatId, username)
		if err != nil {
			log.Fatalf("fail to connect to chat")
		}
	},
}

// Connect to chat
var connectChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "connect to chat with ID",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		// Authorized user and get JWT tokens
		cl.Login(username, password)
		// Save token and endpoint to context
		// Note: only users with admin role can create the chat
		md := metadata.New(map[string]string{"Authorization": authPrefix + cl.accessToken, "Endpoint": connectChatEndpoint})
		ctx = metadata.NewOutgoingContext(ctx, md)
		err := cl.chatClient.Connect(ctx, chatID, username)
		if err != nil {
			log.Fatalf("fail to connect to chat")
		}
	},
}

// Create user command
var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Use to register new user in database",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		id, err := cl.authClient.Create(ctx, &model.User{
			Name:            username,
			Password:        password,
			PasswordConfirm: passwordConfirm,
			Role:            role,
		})
		if err != nil {
			log.Printf("user creation error: %s", err.Error())
		}
		log.Printf("user %s created with id: %d", username, id)
	},
}

// Delete user command
var deleteUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Use to delete user account from database",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		err := cl.authClient.Delete(ctx, userID)
		if err != nil {
			log.Printf("fail to delete user: %s", err.Error())
		}
		log.Printf("user id: %d deleted\n", userID)
	},
}

// Update user command
var updateUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Use to update user account in database",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		err := cl.authClient.Update(ctx, &model.User{
			ID:    userID,
			Name:  username,
			Email: email,
			Role:  role,
		})
		if err != nil {
			log.Printf("fail to update user: %s", err.Error())
		}
		log.Printf("user id: %d updated\n", userID)
	},
}

// Get user command
var getUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Use to get user info",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		user, err := cl.authClient.Get(ctx, userID)
		if err != nil {
			log.Printf("fail to get user: %s", err.Error())
		}
		log.Printf("\nUserInfo:\n name: %s\n role: %s\n email: %s created: %s\n",
			user.Name, user.Role, user.Email, user.CreatedAt)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var err error

	cl, err = NewClientProvider()
	if err != nil {
		os.Exit(1)
	}

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(connectCmd)

	createCmd.AddCommand(createUserCmd)
	getCmd.AddCommand(getUserCmd)
	updateCmd.AddCommand(updateUserCmd)
	deleteCmd.AddCommand(deleteUserCmd)

	createCmd.AddCommand(createChatCmd)
	connectCmd.AddCommand(connectChatCmd)

	// User API
	createUserCmd.Flags().StringVarP(&username, "username", "u", "", "Имя пользователя")
	err := createUserCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark username flag as required: %s\n", err.Error())
	}

	createUserCmd.Flags().StringVarP(&password, "password", "p", "", "User password")
	err = createUserCmd.MarkFlagRequired("password")
	if err != nil {
		log.Fatalf("failed to mark password flag as required: %s\n", err.Error())
	}

	createUserCmd.Flags().StringVarP(&passwordConfirm, "password-confirm", "c", "", "Confirm user password")
	err = createUserCmd.MarkFlagRequired("password-confirm")
	if err != nil {
		log.Fatalf("failed to mark password flag as required: %s\n", err.Error())
	}

	createUserCmd.Flags().StringVarP(&role, "role", "r", "", "User role - choose USER or ADMIN")

	updateUserCmd.Flags().Int64VarP(&userID, "id", "i", 0, "User ID")
	err = updateUserCmd.MarkFlagRequired("id")
	if err != nil {
		log.Fatalf("failed to mark username flag as required: %s\n", err.Error())
	}
	updateUserCmd.Flags().StringVarP(&username, "username", "u", "", "username")
	updateUserCmd.Flags().StringVarP(&email, "email", "e", "", "email")
	updateUserCmd.Flags().StringVarP(&role, "role", "r", "", "role")

	deleteUserCmd.Flags().Int64VarP(&userID, "id", "i", 0, "User ID")
	err = deleteUserCmd.MarkFlagRequired("id")
	if err != nil {
		log.Fatalf("failed to mark username flag as required: %s\n", err.Error())
	}

	getUserCmd.Flags().Int64VarP(&userID, "id", "i", 0, "User ID")
	err = getUserCmd.MarkFlagRequired("id")
	if err != nil {
		log.Fatalf("failed to mark username flag as required: %s\n", err.Error())
	}

	// Create chat command flags
	createChatCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	err = createChatCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark username flag as required: %s\n", err.Error())
	}

	createChatCmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	err = createChatCmd.MarkFlagRequired("password")
	if err != nil {
		log.Fatalf("failed to mark password flag as required: %s\n", err.Error())
	}

	// Connect command flags
	connectChatCmd.Flags().StringVarP(&chatID, "chat-id", "c", "", "Chat ID")
	err = connectChatCmd.MarkFlagRequired("chat-id")
	if err != nil {
		log.Fatalf("failed to mark chat-id flag as required: %s\n", err.Error())
	}

	connectChatCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	err = connectChatCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark username flag as required: %s\n", err.Error())
	}

	connectChatCmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	err = connectChatCmd.MarkFlagRequired("password")
	if err != nil {
		log.Fatalf("failed to mark password flag as required: %s\n", err.Error())
	}
}
