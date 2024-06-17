package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageConnector_FindShortenedURL(t *testing.T) {
	type want struct {
		longURL []byte
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
			stC: &Repository{MapURL: map[string][]byte{
				"aHR0cHM6Ly95YW5kZXgucnU": []byte("https://yandex.ru")}},
			shortenedURL: "aHR0cHM6Ly95YW5kZXgucnU",
			want: want{
				longURL: []byte("https://yandex.ru"),
				found:   true,
			},
		},
		{
			name:         "Negative test #2",
			stC:          &Repository{MapURL: map[string][]byte{}},
			shortenedURL: "aHR0cHM6Ly95YW5kZXgucnU",
			want: want{
				longURL: []byte(nil),
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
		longURL      []byte
	}
	tests := []struct {
		name         string
		stC          *Repository
		longURL      []byte
		shortenedURL string
		want         want
	}{
		{
			name:         "Positive test #1",
			stC:          &Repository{MapURL: map[string][]byte{}},
			longURL:      []byte("https://yandex.ru"),
			shortenedURL: "aHR0cHM6Ly95YW5kZXgucnU",
			want: want{
				shortenedURL: "aHR0cHM6Ly95YW5kZXgucnU",
				longURL:      []byte("https://yandex.ru"),
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
