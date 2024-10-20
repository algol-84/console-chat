package converter

import (
	modelRepo "github.com/algol-84/auth/internal/repository/auth/model"
	desc "github.com/algol-84/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToUserFromRepo конвертирует из модели репо-слоя в тип данных протобафа
func ToUserFromRepo(user *modelRepo.User) *desc.UserInfo {
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

func convertStringToRole(roleStr string) desc.Role {
	if roleValue, exists := desc.Role_value[roleStr]; exists {
		return desc.Role(roleValue)
	}
	return desc.Role_UNKNOWN
}
