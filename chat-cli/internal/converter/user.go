package converter

import (
	"github.com/algol-84/chat-cli/internal/model"
	desc "github.com/algol-84/chat-cli/pkg/user_v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ToUserInfoFromService конвертирует тип model.User в protobuf desc.UserInfo
func ToUserFromModel(user *model.User) *desc.User {
	return &desc.User{
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
		Role:            StringToRole(user.Role),
	}
}

// ToUserInfoFromService конвертирует тип model.User в protobuf desc.UserInfo
func ToModelFromUser(user *desc.UserInfo) *model.User {
	return &model.User{
		ID:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role.String(),
		CreatedAt: user.CreatedAt.AsTime(),
		UpdatedAt: user.UpdatedAt.AsTime(),
	}
}

func FromModelToUserUpdate(user *model.User) *desc.UserUpdate {
	return &desc.UserUpdate{
		Id:    user.ID,
		Name:  wrapperspb.String(user.Name),
		Email: wrapperspb.String(user.Email),
		Role:  StringToRole(user.Role),
	}
}

// StringToRole конвертирует строку в тип desc.Role
func StringToRole(roleStr string) desc.Role {
	if roleValue, exists := desc.Role_value[roleStr]; exists {
		return desc.Role(roleValue)
	}
	return desc.Role_UNKNOWN
}
