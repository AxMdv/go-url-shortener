package service

import "github.com/AxMdv/go-url-shortener/internal/storage"

type shortenerService struct {
	urlRepository storage.Repository
}

// NewShortenerService returns new shortenerService.
func NewShortenerService(urlRepository storage.Repository) *shortenerService {
	// storage.NewRepository()
	return &shortenerService{
		urlRepository: urlRepository,
	}
}
