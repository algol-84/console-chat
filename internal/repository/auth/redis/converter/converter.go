package converter

import (
	"database/sql"
	"time"

	"github.com/algol-84/auth/internal/model"
	modelRepo "github.com/algol-84/auth/internal/repository/auth/redis/model"
)

// ToUserFromRepo конвертирует из модели репо слоя в модель редиса
func ToUserFromRepo(user *modelRepo.User) *model.User {
	var updatedAt sql.NullTime
	if user.UpdatedAtNs != nil {
		updatedAt = sql.NullTime{
			Time:  time.Unix(0, *user.UpdatedAtNs),
			Valid: true,
		}
	}

	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: time.Unix(0, user.CreatedAtNs),
		UpdatedAt: updatedAt,
	}
}

// ToRepoFromUser конвертирует из модели редис в репо-модель
func ToRepoFromUser(user *model.User) *modelRepo.User {
	var updatedAtNs int64
	if user.UpdatedAt.Valid {
		updatedAtNs = user.UpdatedAt.Time.Unix()
	}

	return &modelRepo.User{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Role:        user.Role,
		CreatedAtNs: user.CreatedAt.Unix(),
		UpdatedAtNs: &updatedAtNs,
	}
}
