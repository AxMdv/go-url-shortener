package handlers

import (
	"encoding/base64"

	"github.com/AxMdv/go-url-shortener/internal/storage"
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
func (rb *RequestBatch) ToFormed(uuid string) []storage.FormedURL {
	countReqBatch := len(rb.BatchList)
	urlData := make([]storage.FormedURL, countReqBatch)

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

// ........................................................
