// Package service is service layer of shortener app.
package handlers

import (
	"github.com/AxMdv/go-url-shortener/internal/model"
)

// IShortenerService is the interface that manages shortening urls.
type IShortenerService interface {
	CreateShortURL(formedURL *model.FormedURL) error
	GetLongURL(string) (string, error)
	ShortenLongURL([]byte) string
	CreateShortURLBatch([]model.FormedURL) error
	PingDatabase() error
	GetAllURLByID(string) ([]model.FormedURL, error)
	DeleteURLBatch(model.DeleteBatch) error
	GetFlagByShortURL(shortURL string) (isDeleted bool, err error)
}
