package storage

type RamRepository struct {
	MapURL map[string]string
}

func NewRAMRepository() (*RamRepository, error) {
	return &RamRepository{MapURL: make(map[string]string)}, nil
}

func (rr *RamRepository) AddURL(formedURL *FormedURL) error {
	rr.MapURL[formedURL.ShortenedURL] = formedURL.LongURL
	return nil
}

func (rr *RamRepository) GetURL(shortenedURL string) (string, bool) {
	longURL := rr.MapURL[shortenedURL]
	if longURL == "" {
		return "", false
	}
	return longURL, true
}

func (rr *RamRepository) Close() error {
	return nil
}
