package url

import "github.com/AxMdv/go-url-shortener/internal/storage"

type service struct {
	urlRepository storage.Repository
}

func NewShortenerService(urlRepository storage.Repository) *service {
	// storage.NewRepository()
	return &service{
		urlRepository: urlRepository,
	}
}
