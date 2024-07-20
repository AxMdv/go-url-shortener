package storage

import (
	"github.com/AxMdv/go-url-shortener/internal/app/config"
)

type Repository interface {
	AddURL(*FormedURL) error
	AddURLBatch([]FormedURL) error
	GetURL(shortenedURL string) (string, bool)
	Close() error
}

func NewRepository(config *config.Options) (Repository, error) {
	if config.DataBaseDSN != "" {
		return NewDBRepository(config)
	}
	if config.FileStorage != "" {
		return NewFileRepository(config)
	}
	return NewRAMRepository()
}

type FormedURL struct {
	UIID         string `json:"uiid"`
	ShortenedURL string `json:"short_url"`
	LongURL      string `json:"original_url"`
}
