package service

import (
	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/AxMdv/go-url-shortener/internal/storage"
)

type ShortenerService interface {
	CreateShortURL(*storage.FormedURL) error
	GetLongURL(string) (string, error)
	ShortenLongURL([]byte) string
	CreateShortURLBatch([]storage.FormedURL) error
	PingDatabase(*config.Options) error
	GetAllURLById(string) ([]storage.FormedURL, error)
}
