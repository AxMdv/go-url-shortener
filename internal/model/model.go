// Package model stores global data types models.
package model

// FormedURL is data of an url.
type FormedURL struct {
	UUID          string `json:"uuid,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
	ShortenedURL  string `json:"short_url,omitempty"`
	LongURL       string `json:"original_url,omitempty"`
	DeletedFlag   bool   `json:"deleted_flag,omitempty"`
}

// DeleteBatch is struct to delete batch of urls.
type DeleteBatch struct {
	ShortenedURL []string
	UUID         string
}
