package storage

type StorageConnector struct {
	MapURL map[string][]byte
}

func (stC *StorageConnector) AddURL(longURL []byte, shortenedURL string) {
	stC.MapURL[shortenedURL] = longURL
}

func (stC *StorageConnector) FindShortenedURL(shortenedURL string) ([]byte, bool) {
	longURL := stC.MapURL[shortenedURL]
	found := true
	if longURL == nil {
		found = false
	}
	return longURL, found
}
