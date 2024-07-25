package url

import (
	"fmt"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/AxMdv/go-url-shortener/internal/storage"
)

func (s *service) PingDatabase(config *config.Options) error {
	value, ok := s.urlRepository.(storage.Pinger)

	if !ok {
		return fmt.Errorf("current repo is not database repo")
	}
	err := value.PingDatabase(*config)
	if err != nil {
		return err
	}
	return nil
}
