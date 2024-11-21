// Package service is service layer package of an app.
package service

type ShortenerService struct {
	urlRepository IRepository
}

// NewShortenerService returns new shortenerService.
func NewShortenerService(urlRepository IRepository) *ShortenerService {
	// storage.NewRepository()
	return &ShortenerService{
		urlRepository: urlRepository,
	}
}
