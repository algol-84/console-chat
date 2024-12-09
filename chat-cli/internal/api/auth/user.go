package auth

import (
	"context"
	"log"

	"github.com/algol-84/chat-cli/internal/converter"
	"github.com/algol-84/chat-cli/internal/model"
	desc "github.com/algol-84/chat-cli/pkg/user_v1"
)

const (
	authHost = "127.0.0.1"
	authPort = 50051
)

var UserClient desc.UserV1Client

func (a *AuthImpl) Create(user *model.User) (int64, error) {
	ctx := context.Background()
	res, err := a.userClient.Create(ctx, &desc.CreateRequest{
		User: converter.ToUserFromModel(user),
	})

	if err != nil {
		return 0, err
	}

	return res.Id, nil
}

func (a *AuthImpl) Get(userID int64) (*model.User, error) {
	ctx := context.Background()
	user, err := a.userClient.Get(ctx, &desc.GetRequest{
		Id: userID,
	})
	if err != nil {
		return nil, err
	}

	return converter.ToModelFromUser(user.UserInfo), nil
}

func (a *AuthImpl) Delete(userID int64) error {
	ctx := context.Background()
	_, err := a.userClient.Delete(ctx, &desc.DeleteRequest{
		Id: userID,
	})

	if err != nil {
		return err
	}

	return nil
}

func (a *AuthImpl) Update(user *model.User) error {
	ctx := context.Background()
	log.Println("update handler")
	_, err := a.userClient.Update(ctx, &desc.UpdateRequest{
		UserUpdate: converter.FromModelToUserUpdate(user),
	})

	if err != nil {
		return err
	}

	return nil
}
