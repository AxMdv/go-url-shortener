package service

import (
	"testing"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/AxMdv/go-url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewShortenerService(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        "",
	}
	repository, err := storage.NewRepository(config)
	require.NoError(t, err)
	_ = NewShortenerService(repository)
}

func TestServiceCreateShortURL(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        "",
	}
	repository, err := storage.NewRepository(config)
	require.NoError(t, err)
	shortenerService := NewShortenerService(repository)

	type want struct {
		err error
	}
	tests := []struct {
		name      string
		formedURL *storage.FormedURL
		want      want
	}{
		{
			name: "Positive test #1",
			formedURL: &storage.FormedURL{
				LongURL:      "https://vk.com",
				ShortenedURL: "http://localhost:8080/aHR0cHM6Ly92ay5jb20",
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		err = shortenerService.CreateShortURL(tt.formedURL)
		require.NoError(t, err)
	}
}

func TestServiceCreateShortURLBatch(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        "",
	}
	repository, err := storage.NewRepository(config)
	require.NoError(t, err)
	shortenerService := NewShortenerService(repository)

	type want struct {
		err error
	}
	tests := []struct {
		name      string
		formedURL []storage.FormedURL
		want      want
	}{
		{
			name: "Positive test #1",
			formedURL: []storage.FormedURL{
				{
					LongURL:      "https://vk.com",
					ShortenedURL: "http://localhost:8080/aHR0cHM6Ly92ay5jb20",
				},
				{
					LongURL:      "https://yandex.ru",
					ShortenedURL: "http://localhost:8080/aHR0cHM6Ly95YW5kZXgucnU",
				},
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		err = shortenerService.CreateShortURLBatch(tt.formedURL)
		require.NoError(t, err)
	}
}

func TestServiceShortenLongURL(t *testing.T) {
	service := shortenerService{}
	type want struct {
		shortURL string
	}
	tests := []struct {
		name    string
		longURL string
		want    want
	}{
		{
			name:    "Positive test #1",
			longURL: "https://vk.com",

			want: want{
				shortURL: "aHR0cHM6Ly92ay5jb20",
			},
		},
		{
			name:    "Positive test #2",
			longURL: "https://yandex.ru",
			want: want{
				shortURL: "aHR0cHM6Ly95YW5kZXgucnU",
			},
		},
	}
	for _, tt := range tests {
		shortURL := service.ShortenLongURL([]byte(tt.longURL))
		assert.Equal(t, tt.want.shortURL, shortURL)
	}
}

func TestServiceDeleteURLBatch(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        "",
	}
	repository, err := storage.NewRepository(config)
	require.NoError(t, err)
	shortenerService := NewShortenerService(repository)

	type want struct {
		err error
	}
	tests := []struct {
		name        string
		deleteBatch storage.DeleteBatch
		want        want
	}{
		{
			name: "Positive test #1",
			deleteBatch: storage.DeleteBatch{
				ShortenedURL: []string{"aHR0cHM6Ly92ay5jb20", "aHR0cHM6Ly95YW5kZXgucnU"},
				UUID:         "asd",
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		err = shortenerService.DeleteURLBatch(tt.deleteBatch)
		require.NoError(t, err)
	}
}

func TestServiceGetLongURL(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        "",
	}
	repository, err := storage.NewRepository(config)
	require.NoError(t, err)
	shortenerService := NewShortenerService(repository)

	type want struct {
		longURL string
		err     error
	}
	tests := []struct {
		name      string
		formedURL *storage.FormedURL
		want      want
	}{
		{
			name: "Positive test #1",
			formedURL: &storage.FormedURL{
				LongURL:      "https://vk.com",
				ShortenedURL: "aHR0cHM6Ly92ay5jb20",
			},
			want: want{
				longURL: "https://vk.com",
				err:     nil,
			},
		},
	}
	for _, tt := range tests {
		err = shortenerService.CreateShortURL(tt.formedURL)
		require.NoError(t, err)
		longURL, err := shortenerService.GetLongURL(tt.formedURL.ShortenedURL)
		assert.Equal(t, tt.want.err, err)
		assert.Equal(t, tt.want.longURL, longURL)
	}
}

func TestServiceGetAllURLByID(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        "",
	}
	repository, err := storage.NewRepository(config)
	require.NoError(t, err)
	shortenerService := NewShortenerService(repository)

	type want struct {
		formedURL []storage.FormedURL
		err       error
	}
	tests := []struct {
		name      string
		formedURL []storage.FormedURL
		uuid      string
		want      want
	}{
		{
			name: "Positive test #1",
			formedURL: []storage.FormedURL{
				{
					LongURL:      "https://vk.com",
					ShortenedURL: "aHR0cHM6Ly92ay5jb20",
					UUID:         "asd",
				},
			},
			uuid: "asd",
			want: want{
				formedURL: []storage.FormedURL{
					{
						LongURL:      "https://vk.com",
						ShortenedURL: "aHR0cHM6Ly92ay5jb20",
					},
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		// first add urls to db
		for _, formedURL := range tt.formedURL {
			err = shortenerService.CreateShortURL(&formedURL)
			require.NoError(t, err)
		}

		formedURL, err := shortenerService.GetAllURLByID(tt.uuid)
		assert.Equal(t, tt.want.err, err)
		assert.Equal(t, tt.want.formedURL, formedURL)
	}
}

func TestServiceGetFlagByShortURL(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        "",
	}
	repository, err := storage.NewRepository(config)
	require.NoError(t, err)
	shortenerService := NewShortenerService(repository)

	type want struct {
		isDeleted bool
		err       error
	}
	tests := []struct {
		name      string
		formedURL storage.FormedURL
		uuid      string
		want      want
	}{
		{
			name: "Positive test #1",
			formedURL: storage.FormedURL{
				LongURL:      "https://vk.com",
				ShortenedURL: "aHR0cHM6Ly92ay5jb20",
				UUID:         "asd",
			},
			uuid: "asd",
			want: want{
				isDeleted: true,
				err:       nil,
			},
		},
	}
	for _, tt := range tests {
		// add url
		err := shortenerService.CreateShortURL(&tt.formedURL)
		require.NoError(t, err)
		// delete url
		err = shortenerService.DeleteURLBatch(storage.DeleteBatch{UUID: tt.uuid, ShortenedURL: []string{tt.formedURL.ShortenedURL}})
		require.NoError(t, err)

		formedURL, err := shortenerService.GetFlagByShortURL(tt.formedURL.ShortenedURL)
		assert.Equal(t, tt.want.err, err)
		assert.Equal(t, tt.want.isDeleted, formedURL)
	}
}
