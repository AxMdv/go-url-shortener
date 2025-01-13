package service

import (
	"context"

	"github.com/AxMdv/go-url-shortener/internal/model"
)

// IRepository is the interface that stores urls.
type IRepository interface {
	AddURL(context.Context, *model.FormedURL) error
	AddURLBatch(context.Context, []model.FormedURL) error
	GetURL(context.Context, string) (string, error)
	GetURLByUserID(context.Context, string) ([]model.FormedURL, error)
	DeleteURLBatch(ctx context.Context, formedURL []model.FormedURL) error
	GetFlagByShortURL(context.Context, string) (bool, error)
}
