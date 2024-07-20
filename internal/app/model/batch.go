package model

type RequestBatch struct {
	BatchList []BatchOriginal
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
