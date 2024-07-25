package url

import (
	"encoding/base64"

	"github.com/AxMdv/go-url-shortener/internal/storage"
)

func (s *service) CreateShortURL(formedURL *storage.FormedURL) error {
	err := s.urlRepository.AddURL(formedURL)
	return err
}

func (s *service) CreateShortURLBatch(formedURL []storage.FormedURL) error {
	err := s.urlRepository.AddURLBatch(formedURL)
	return err
}

//.......................................

func (s *service) ShortenLongURL(longURL []byte) string {
	shortenedURL := base64.RawStdEncoding.EncodeToString(longURL)
	return shortenedURL
}
