package storage

import (
	"errors"
	"fmt"

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

// .............................................................

type AddURLError struct {
	DuplicateValue string
	Err            error
}

func (ae *AddURLError) Error() string {
	return fmt.Sprintf("%v %v", ae.DuplicateValue, ae.Err)
}
func NewDuplicateError(err error, shortenedURL string) error {
	return &AddURLError{
		DuplicateValue: shortenedURL,
		Err:            err,
	}
}

var ErrDuplicate = errors.New("url already exists")

// .............................................................
