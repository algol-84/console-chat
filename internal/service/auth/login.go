package auth

import (
	"context"
	"log"
)

func (s *service) Login(_ context.Context) error {
	log.Printf("login handle in service layer")
	return nil
}
