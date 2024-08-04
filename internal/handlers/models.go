package handlers

import (
	"encoding/base64"

	"github.com/AxMdv/go-url-shortener/internal/storage"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

// ........................................................
type RequestBatch struct {
	BatchList []BatchOriginal
}

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

type BatchOriginal struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResponseBatch struct {
	BatchList []BatchShortened
}

type BatchShortened struct {
	CorrelationID string `json:"correlation_id"`
	ShortenedURL  string `json:"short_url"`
}

// ........................................................
