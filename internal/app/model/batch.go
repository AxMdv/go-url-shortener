package model

type RequestBatch struct {
	BatchList []BatchOriginal
}

type BatchOriginal struct {
	CorrelationId string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResponseBatch struct {
	BatchList []BatchShortened
}

type BatchShortened struct {
	CorrelationId string `json:"correlation_id"`
	ShortenedURL  string `json:"short_url"`
}
