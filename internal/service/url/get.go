package url

import (
	"context"
	"time"

	"github.com/AxMdv/go-url-shortener/internal/storage"
)

func (s *service) GetLongURL(shortenedURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	longURL, err := s.urlRepository.GetURL(ctx, shortenedURL)
	if err != nil {
		return "", err
	}
	return longURL, nil
}

func (s *service) GetAllURLByID(uuid string) ([]storage.FormedURL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	formedURL, err := s.urlRepository.GetURLByUserID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return formedURL, nil

}

func (s *service) GetFlagByShortURL(shortURL string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	isDeleted, err := s.urlRepository.GetFlagByShortURL(ctx, shortURL)

	return isDeleted, err
}
