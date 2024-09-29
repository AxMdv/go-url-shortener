package service

import (
	"github.com/AxMdv/go-url-shortener/internal/storage"
)

// ShortenerService is the interface that manages shortening urls.
type ShortenerService interface {
	CreateShortURL(*storage.FormedURL) error
	GetLongURL(string) (string, error)
	ShortenLongURL([]byte) string
	CreateShortURLBatch([]storage.FormedURL) error
	PingDatabase() error
	GetAllURLByID(string) ([]storage.FormedURL, error)
	DeleteURLBatch(storage.DeleteBatch)
	GetFlagByShortURL(shortURL string) (isDeleted bool, err error)
}
