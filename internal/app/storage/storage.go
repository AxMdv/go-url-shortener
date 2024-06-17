package storage

type Repository struct {
	MapURL map[string][]byte
}

func (r *Repository) AddURL(longURL []byte, shortenedURL string) {
	r.MapURL[shortenedURL] = longURL
}

func (stC *Repository) FindShortenedURL(shortenedURL string) ([]byte, bool) {
	longURL := stC.MapURL[shortenedURL]
	found := true
	if longURL == nil {
		found = false
	}
	return longURL, found
}
