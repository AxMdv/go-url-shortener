package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/AxMdv/go-url-shortener/internal/config"
)

type Repository interface {
	AddURL(context.Context, *FormedURL) error
	AddURLBatch(context.Context, []FormedURL) error
	GetURL(context.Context, string) (string, error)
	GetURLByUserID(context.Context, string) ([]FormedURL, error)
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
	UUID          string `json:"uuid,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
	ShortenedURL  string `json:"short_url,omitempty"`
	LongURL       string `json:"original_url,omitempty"`
}

// func (fu []FormedURL) ToResponseBatch() handlers.ResponseBatch {
// 	respData := make([]handlers.BatchShortened, len(fu))
// 	for i, v := range fu {
// 		respData[i].CorrelationID = v.UIID
// 		respData[i].ShortenedURL = fmt.Sprintf("%s/%s", s.Config.ResponseResultAddr, v.ShortenedURL)
// 	}
// 	return respData
// }

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

type NoContentError struct {
	UIID string
	Err  error
}

func (nc *NoContentError) Error() string {
	return fmt.Sprintf("%v %v", nc.UIID, nc.Err)
}

func NewNoContentError(err error, uuid string) error {
	return &NoContentError{
		UIID: uuid,
		Err:  err,
	}
}

var ErrNoContent = errors.New("no urls created by current user ")

// .............................................................
type Pinger interface {
	PingDB(context.Context, config.Options) error
}

type Closer interface {
	Close() error
}
