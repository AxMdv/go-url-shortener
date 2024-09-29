package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/AxMdv/go-url-shortener/internal/config"
)

// Repository is the interface that stores urls.
type Repository interface {
	AddURL(context.Context, *FormedURL) error
	AddURLBatch(context.Context, []FormedURL) error
	GetURL(context.Context, string) (string, error)
	GetURLByUserID(context.Context, string) ([]FormedURL, error)
	DeleteURLBatch(ctx context.Context, formedURL []FormedURL) error
	GetFlagByShortURL(context.Context, string) (bool, error)
}

// NewRepository returns repository, that implements Repository interface.
func NewRepository(config *config.Options) (Repository, error) {
	if config.DataBaseDSN != "" {
		return NewDBRepository(config)
	}
	if config.FileStorage != "" {
		return NewFileRepository(config)
	}
	return NewRAMRepository()
}

// FormedURL is data of an url.
type FormedURL struct {
	UUID          string `json:"uuid,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
	ShortenedURL  string `json:"short_url,omitempty"`
	LongURL       string `json:"original_url,omitempty"`
	DeletedFlag   bool   `json:"deleted_flag,omitempty"`
}

// .............................................................

// AddURLError is duplicating error.
type AddURLError struct {
	DuplicateValue string
	Err            error
}

// Error is provided to implement Error interface.
func (ae *AddURLError) Error() string {
	return fmt.Sprintf("%v %v", ae.DuplicateValue, ae.Err)
}

// NewDuplicateError return new AddURLError.
func NewDuplicateError(err error, shortenedURL string) error {
	return &AddURLError{
		DuplicateValue: shortenedURL,
		Err:            err,
	}
}

// ErrDuplicate is duplicate error.
var ErrDuplicate = errors.New("url already exists")

// .............................................................

// NoContentError is error when there is no wanted data.
type NoContentError struct {
	UIID string
	Err  error
}

// Error is provided to implement Error interface.
func (nc *NoContentError) Error() string {
	return fmt.Sprintf("%v %v", nc.UIID, nc.Err)
}

// NewNoContentError returns new NoContentError.
func NewNoContentError(err error, uuid string) error {
	return &NoContentError{
		UIID: uuid,
		Err:  err,
	}
}

// ErrNoContent is  error when there is no wanted data.
var ErrNoContent = errors.New("no urls created by current user ")

// .............................................................

// Pinger is an interface that pings DB.
type Pinger interface {
	PingDB(context.Context) error
}

// Closer is an interface that can close DB.
type Closer interface {
	Close() error
}

// .............................................................

// DeleteBatch is stuct to delete batch of urls.
type DeleteBatch struct {
	ShortenedURL []string
	UUID         string
}
