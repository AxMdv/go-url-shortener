package service

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/AxMdv/go-url-shortener/internal/model"
)

// CreateShortURL shortens long url.
func (s *ShortenerService) CreateShortURL(formedURL *model.FormedURL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := s.urlRepository.AddURL(ctx, formedURL)
	return err
}

// CreateShortURLBatch shortens batch of long urls.
func (s *ShortenerService) CreateShortURLBatch(formedURL []model.FormedURL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := s.urlRepository.AddURLBatch(ctx, formedURL)
	return err
}

// ShortenLongURL transform long url to short url.
func (s *ShortenerService) ShortenLongURL(longURL []byte) string {
	shortenedURL := base64.RawStdEncoding.EncodeToString(longURL)
	return shortenedURL
}
