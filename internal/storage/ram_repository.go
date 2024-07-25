package storage

import "context"

type RAMRepository struct {
	MapURL map[string]string
}

func NewRAMRepository() (*RAMRepository, error) {
	return &RAMRepository{MapURL: make(map[string]string)}, nil
}

func (rr *RAMRepository) AddURL(_ context.Context, formedURL *FormedURL) error {
	if rr.MapURL[formedURL.ShortenedURL] != "" {
		return NewDuplicateError(ErrDuplicate, formedURL.ShortenedURL)
	}
	rr.MapURL[formedURL.ShortenedURL] = formedURL.LongURL
	return nil
}

func (rr *RAMRepository) AddURLBatch(_ context.Context, formedURL []FormedURL) error {
	for _, v := range formedURL {
		rr.MapURL[v.ShortenedURL] = v.LongURL
	}
	return nil
}

func (rr *RAMRepository) GetURL(_ context.Context, shortenedURL string) (string, error) {
	longURL := rr.MapURL[shortenedURL]
	if longURL == "" {
		return "", nil
	}
	return longURL, nil
}
