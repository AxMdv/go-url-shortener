package storage

import "context"

type RAMRepository struct {
	MapURL     map[string]string   //[shortened]long
	MapUUID    map[string][]string //[uuid][]shortened
	MapDeleted map[string]bool     //[shortened]deleted_flag
}

func NewRAMRepository() (*RAMRepository, error) {
	return &RAMRepository{MapURL: make(map[string]string), MapUUID: make(map[string][]string)}, nil
}

func (rr *RAMRepository) AddURL(_ context.Context, formedURL *FormedURL) error {
	if rr.MapURL[formedURL.ShortenedURL] != "" {
		return NewDuplicateError(ErrDuplicate, formedURL.ShortenedURL)
	}
	rr.MapURL[formedURL.ShortenedURL] = formedURL.LongURL
	rr.MapUUID[formedURL.UUID] = append(rr.MapUUID[formedURL.UUID], formedURL.ShortenedURL)
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

func (rr *RAMRepository) GetURLByUserID(_ context.Context, uuid string) ([]FormedURL, error) {
	shortenedURL := rr.MapUUID[uuid]
	formedURL := make([]FormedURL, 0)
	for _, v := range shortenedURL {
		longURL, err := rr.GetURL(context.Background(), v)
		if err != nil {
			return nil, err
		}
		var fu FormedURL
		fu.LongURL = longURL
		fu.ShortenedURL = v
		formedURL = append(formedURL, fu)

	}
	return formedURL, nil
}

func (rr *RAMRepository) DeleteURLBatch(ctx context.Context, formedURL []FormedURL) error {

	for _, v := range formedURL {

		sliceShortened := rr.MapUUID[v.UUID]
		if sliceShortened == nil {
			continue
		}

		contains := contains(sliceShortened, v.ShortenedURL)
		if !contains {
			continue
		}
		rr.MapDeleted[v.ShortenedURL] = true
	}
	return nil
}

func contains(target []string, value string) bool {
	for _, v := range target {
		if v == value {
			return true
		}
	}
	return false
}

func (rr *RAMRepository) GetFlagByShortURL(_ context.Context, shortenedURL string) (bool, error) {
	return rr.MapDeleted[shortenedURL], nil
}
