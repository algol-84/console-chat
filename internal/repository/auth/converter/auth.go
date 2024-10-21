package converter

import (
	model "github.com/algol-84/auth/internal/model"
	modelRepo "github.com/algol-84/auth/internal/repository/auth/model"
)

// ToUserFromRepo конвертирует из модели репо слоя в модель сервисного слоя
func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Password:  user.Password,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
