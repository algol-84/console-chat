package auth

import (
	"context"
	"log"
)

// Check определяет уровень доступа пользователя
func (s *service) Check(_ context.Context) error {
	log.Printf("check handle in service layer")
	return nil
}
