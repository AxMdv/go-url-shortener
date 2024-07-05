package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageConnector_FindShortenedURL(t *testing.T) {
	type want struct {
		longURL string
		found   bool
	}
	tests := []struct {
		name         string
		stC          *Repository
		shortenedURL string
		want         want
	}{
		{
			name: "Positive test #1",
			stC: &Repository{MapURL: map[string]string{
				"aHR0cHM6Ly95YW5kZXgucnU": "https://yandex.ru"}},
			shortenedURL: "aHR0cHM6Ly95YW5kZXgucnU",
			want: want{
				longURL: "https://yandex.ru",
				found:   true,
			},
		},
		{
			name:         "Negative test #2",
			stC:          &Repository{MapURL: map[string]string{}},
			shortenedURL: "aHR0cHM6Ly95YW5kZXgucnU",
			want: want{
				longURL: "",
				found:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			longURL, found := tt.stC.FindShortenedURL(tt.shortenedURL)
			assert.Equal(t, tt.want.longURL, longURL)
			assert.Equal(t, tt.want.found, found)
		})
	}
}

func TestStorageConnector_AddURL(t *testing.T) {
	type want struct {
		shortenedURL string
		longURL      string
	}
	tests := []struct {
		name         string
		stC          *Repository
		longURL      string
		shortenedURL string
		want         want
	}{
		{
			name:         "Positive test #1",
			stC:          &Repository{MapURL: map[string]string{}},
			longURL:      "https://yandex.ru",
			shortenedURL: "aHR0cHM6Ly95YW5kZXgucnU",
			want: want{
				shortenedURL: "aHR0cHM6Ly95YW5kZXgucnU",
				longURL:      "https://yandex.ru",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stC.AddURL(tt.longURL, tt.shortenedURL)
			assert.Equal(t, tt.want.longURL, tt.stC.MapURL[tt.shortenedURL])
		})
	}
}
