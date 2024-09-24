package handlers_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/AxMdv/go-url-shortener/internal/handlers"
	"github.com/AxMdv/go-url-shortener/internal/shortener"
	"github.com/AxMdv/go-url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateShortURL(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        "",
	}
	repository, err := storage.NewRepository(config)
	require.NoError(t, err)
	urlService := shortener.NewShortenerService(repository)
	shortenerHandlers := handlers.NewShortenerHandlers(urlService, config)

	type want struct {
		contentType string
		respBody    string
		statusCode  int
	}
	tests := []struct {
		name       string
		requestURL string
		reqBody    string
		want       want
	}{
		{
			name:       "Positive test #1",
			requestURL: "/",
			reqBody:    "https://yandex.ru",
			want: want{
				contentType: "text/plain",
				respBody:    "http://localhost:8080/aHR0cHM6Ly95YW5kZXgucnU",
				statusCode:  201,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			reqBody := bytes.NewReader([]byte(tt.reqBody))
			request := httptest.NewRequest(http.MethodPost, tt.requestURL, reqBody)
			w := httptest.NewRecorder()
			shortenerHandlers.CreateShortURL(w, request)
			result := w.Result()

			resultURL, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			resultString := string(resultURL)

			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.respBody, resultString)
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestShortenerHandlers_CreateShortURLJson(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        "",
	}
	repository, err := storage.NewRepository(config)
	require.NoError(t, err)
	urlService := shortener.NewShortenerService(repository)
	shortenerHandlers := handlers.NewShortenerHandlers(urlService, config)

	type want struct {
		contentType string
		respBody    string
		statusCode  int
	}
	tests := []struct {
		name           string
		requestURL     string
		reqBody        string
		reqContentType string
		want           want
	}{
		{
			name:           "Positive test #1",
			requestURL:     "/api/shorten",
			reqBody:        `{"url": "https://yandex.ru"} `,
			reqContentType: "application/json",
			want: want{
				contentType: "application/json",
				respBody:    `{"result":"http://localhost:8080/aHR0cHM6Ly95YW5kZXgucnU"}`,
				statusCode:  201,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := bytes.NewReader([]byte(tt.reqBody))
			request := httptest.NewRequest(http.MethodPost, tt.requestURL, reqBody)
			request.Header.Add("Content-Type", tt.reqContentType)
			w := httptest.NewRecorder()
			shortenerHandlers.CreateShortURLJson(w, request)
			result := w.Result()

			resultURL, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			resultString := string(resultURL)

			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.respBody, resultString)
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
