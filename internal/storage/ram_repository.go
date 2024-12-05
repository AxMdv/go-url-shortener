package storage

import (
	"context"

	"github.com/AxMdv/go-url-shortener/internal/model"
)

// RAMRepository is in-memory repository.
type RAMRepository struct {
	MapURL     map[string]string   // [shortened]long
	MapUUID    map[string][]string // [uuid][]shortened
	MapDeleted map[string]bool     // [shortened]deleted_flag
}

// NewRAMRepository returns new RAMRepository.
func NewRAMRepository() (*RAMRepository, error) {
	return &RAMRepository{MapURL: make(map[string]string), MapUUID: make(map[string][]string), MapDeleted: make(map[string]bool)}, nil
}

// AddURL writes url to RAMRepository.
func (rr *RAMRepository) AddURL(_ context.Context, formedURL *model.FormedURL) error {
	if rr.MapURL[formedURL.ShortenedURL] != "" {
		return NewDuplicateError(ErrDuplicate, formedURL.ShortenedURL)
	}
	rr.MapURL[formedURL.ShortenedURL] = formedURL.LongURL
	rr.MapUUID[formedURL.UUID] = append(rr.MapUUID[formedURL.UUID], formedURL.ShortenedURL)
	return nil
}

// AddURLBatch writes batch of urls to RAMRepository.
func (rr *RAMRepository) AddURLBatch(_ context.Context, formedURL []model.FormedURL) error {
	for _, v := range formedURL {
		rr.MapURL[v.ShortenedURL] = v.LongURL
	}
	return nil
}

// GetURL returns url from RAMRepository.
func (rr *RAMRepository) GetURL(_ context.Context, shortenedURL string) (string, error) {
	longURL := rr.MapURL[shortenedURL]
	return longURL, nil
}

// GetURLByUserID returns urls shortened by user from RAMRepository.
func (rr *RAMRepository) GetURLByUserID(_ context.Context, uuid string) ([]model.FormedURL, error) {
	shortenedURL := rr.MapUUID[uuid]
	formedURL := make([]model.FormedURL, 0)
	for _, v := range shortenedURL {
		longURL, err := rr.GetURL(context.Background(), v)
		if err != nil {
			return nil, err
		}
		var fu model.FormedURL
		fu.LongURL = longURL
		fu.ShortenedURL = v
		formedURL = append(formedURL, fu)

	}
	return formedURL, nil
}

// DeleteURLBatch deletes urls created by user.
func (rr *RAMRepository) DeleteURLBatch(ctx context.Context, formedURL []model.FormedURL) error {

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

// GetFlagByShortURL returns if shortened url is deleted.
func (rr *RAMRepository) GetFlagByShortURL(_ context.Context, shortenedURL string) (bool, error) {
	return rr.MapDeleted[shortenedURL], nil
}

func contains(target []string, value string) bool {
	for _, v := range target {
		if v == value {
			return true
		}
	}
	return false
}
