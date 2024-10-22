// Package converter Конвертер типов protobuf <-> service model
package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/algol-84/auth/internal/model"
	desc "github.com/algol-84/auth/pkg/user_v1"
)

// FromUserToService конвертирует тип protobuf User в модель сервисного слоя
func FromUserToService(user *desc.User) *model.User {
	return &model.User{
		Name:            user.Name,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
		Email:           user.Email,
		Role:            user.Role.String(),
	}
}

// ToUserInfoFromService конвертирует тип model.User в protobuf desc.UserInfo
func ToUserInfoFromService(user *model.User) *desc.UserInfo {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.UserInfo{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      convertStringToRole(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// ToUserUpdateFromDesc конвертирует тип protobuf UserUpdate в модель UserUpdate
func ToUserUpdateFromDesc(user *desc.UserUpdate) *model.UserUpdate {
	var userUpdate model.UserUpdate
	if user.Name != nil {
		userUpdate.Name.Value = user.Name.Value
		userUpdate.Name.Valid = true
	}
	if user.Email != nil {
		userUpdate.Email.Value = user.Email.Value
		userUpdate.Email.Valid = true
	}
	if user.Role != desc.Role_UNKNOWN {
		userUpdate.Role.Value = user.Role.String()
		userUpdate.Role.Valid = true
	}
	userUpdate.ID = user.Id

	return &userUpdate
}

func convertStringToRole(roleStr string) desc.Role {
	if roleValue, exists := desc.Role_value[roleStr]; exists {
		return desc.Role(roleValue)
	}
	return desc.Role_UNKNOWN
}
