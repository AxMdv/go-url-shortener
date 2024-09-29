package service

import (
	"context"
	"time"

	"github.com/AxMdv/go-url-shortener/internal/storage"
)

// GetLongURL returns long url if it was shortened earlier.
func (s *shortenerService) GetLongURL(shortenedURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	longURL, err := s.urlRepository.GetURL(ctx, shortenedURL)
	if err != nil {
		return "", err
	}
	return longURL, nil
}

// GetAllURLByID returns all urls shortened by user.
func (s *shortenerService) GetAllURLByID(uuid string) ([]storage.FormedURL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	formedURL, err := s.urlRepository.GetURLByUserID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return formedURL, nil

}

// GetFlagByShortURL returns if short url was deleted.
func (s *shortenerService) GetFlagByShortURL(shortURL string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	isDeleted, err := s.urlRepository.GetFlagByShortURL(ctx, shortURL)

	return isDeleted, err
}
