package chat

import (
	"context"
	"log"
)

func (s *service) Delete(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.chatRepository.Delete(ctx, id)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.logRepository.Create(ctx, id, "chat was deleted")
		if errTx != nil {
			log.Println(errTx)
			return errTx
		}
		return nil
	})

	// if err != nil {
	// 	return nil
	// }

	// err := s.chatRepository.Delete(ctx, id)
	// if err != nil {
	// 	return err
	// }

	return err
}
