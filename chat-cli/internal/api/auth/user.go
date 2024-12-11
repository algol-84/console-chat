package auth

import (
	"context"

	"github.com/algol-84/chat-cli/internal/converter"
	"github.com/algol-84/chat-cli/internal/model"
	desc "github.com/algol-84/chat-cli/pkg/user_v1"
)

var UserClient desc.UserV1Client

func (a *AuthImpl) Create(ctx context.Context, user *model.User) (int64, error) {
	res, err := a.userClient.Create(ctx, &desc.CreateRequest{
		User: converter.ToUserFromModel(user),
	})

	if err != nil {
		return 0, err
	}

	return res.Id, nil
}

func (a *AuthImpl) Get(ctx context.Context, userID int64) (*model.User, error) {
	user, err := a.userClient.Get(ctx, &desc.GetRequest{
		Id: userID,
	})
	if err != nil {
		return nil, err
	}

	return converter.ToModelFromUser(user.UserInfo), nil
}

func (a *AuthImpl) Delete(ctx context.Context, userID int64) error {
	_, err := a.userClient.Delete(ctx, &desc.DeleteRequest{
		Id: userID,
	})

	if err != nil {
		return err
	}

	return nil
}

func (a *AuthImpl) Update(ctx context.Context, user *model.User) error {
	_, err := a.userClient.Update(ctx, &desc.UpdateRequest{
		UserUpdate: converter.FromModelToUserUpdate(user),
	})

	if err != nil {
		return err
	}

	return nil
}
