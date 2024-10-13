package service

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/AxMdv/go-url-shortener/internal/storage"
)

func (s *shortenerService) CreateShortURL(formedURL *storage.FormedURL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := s.urlRepository.AddURL(ctx, formedURL)
	return err
}

func (s *shortenerService) CreateShortURLBatch(formedURL []storage.FormedURL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := s.urlRepository.AddURLBatch(ctx, formedURL)
	return err
}

//.......................................

func (s *shortenerService) ShortenLongURL(longURL []byte) string {
	shortenedURL := base64.RawStdEncoding.EncodeToString(longURL)
	return shortenedURL
}
