package storage

type Repository struct {
	MapURL map[string][]byte
}

func (r *Repository) AddURL(longURL []byte, shortenedURL string) {
	r.MapURL[shortenedURL] = longURL
}

func (r *Repository) FindShortenedURL(shortenedURL string) ([]byte, bool) {
	longURL := r.MapURL[shortenedURL]
	found := true
	if longURL == nil {
		found = false
	}
	return longURL, found
}
