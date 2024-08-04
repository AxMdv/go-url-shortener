package service

import (
	"context"
	"fmt"
	"time"

	"github.com/AxMdv/go-url-shortener/internal/storage"
)

func (s *service) PingDatabase() error {
	value, ok := s.urlRepository.(storage.Pinger)

	if !ok {
		return fmt.Errorf("current repo is not database repo")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := value.PingDB(ctx)
	if err != nil {
		return err
	}
	return nil
}
