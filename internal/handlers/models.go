package handlers

import (
	"encoding/base64"

	"github.com/AxMdv/go-url-shortener/internal/model"
)

// Request is a default request struct.
type Request struct {
	URL string `json:"url"`
}

// Response is a default request struct.
type Response struct {
	Result string `json:"result"`
}

// RequestBatch is a struct of BatchList.
type RequestBatch struct {
	BatchList []BatchOriginal
}

// ToFormed transform RequestBatch to  slice of FormedURL.
func (rb *RequestBatch) ToFormed(uuid string) []model.FormedURL {
	countReqBatch := len(rb.BatchList)
	urlData := make([]model.FormedURL, countReqBatch)

	for i, v := range rb.BatchList {
		urlData[i].UUID = uuid
		urlData[i].CorrelationID = v.CorrelationID
		urlData[i].LongURL = v.OriginalURL
		shortenedURL := base64.RawStdEncoding.EncodeToString([]byte(v.OriginalURL))
		urlData[i].ShortenedURL = shortenedURL
	}
	return urlData
}

// BatchOriginal is requested batch item.
type BatchOriginal struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ResponseBatch is a struct of shortened BatchList.
type ResponseBatch struct {
	BatchList []BatchShortened
}

// BatchShortened is batch url for response.
type BatchShortened struct {
	CorrelationID string `json:"correlation_id"`
	ShortenedURL  string `json:"short_url"`
}

type FormedURL struct {
	UUID          string `json:"uuid,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
	ShortenedURL  string `json:"short_url,omitempty"`
	LongURL       string `json:"original_url,omitempty"`
	DeletedFlag   bool   `json:"deleted_flag,omitempty"`
}

// ........................................................
