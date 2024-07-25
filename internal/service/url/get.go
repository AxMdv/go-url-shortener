package url

import (
	"context"
	"time"
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
