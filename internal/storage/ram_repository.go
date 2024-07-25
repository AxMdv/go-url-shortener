package storage

type RAMRepository struct {
	MapURL map[string]string
}

func NewRAMRepository() (*RAMRepository, error) {
	return &RAMRepository{MapURL: make(map[string]string)}, nil
}

func (rr *RAMRepository) AddURL(formedURL *FormedURL) error {
	if rr.MapURL[formedURL.ShortenedURL] != "" {
		return NewDuplicateError(ErrDuplicate, formedURL.ShortenedURL)
	}
	rr.MapURL[formedURL.ShortenedURL] = formedURL.LongURL
	return nil
}

func (rr *RAMRepository) AddURLBatch(formedURL []FormedURL) error {
	for _, v := range formedURL {
		rr.MapURL[v.ShortenedURL] = v.LongURL
	}
	return nil
}

func (rr *RAMRepository) GetURL(shortenedURL string) (string, error) {
	longURL := rr.MapURL[shortenedURL]
	if longURL == "" {
		return "", nil
	}
	return longURL, nil
}

func (rr *RAMRepository) Close() error {
	return nil
}
