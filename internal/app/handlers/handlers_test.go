package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AxMdv/go-url-shortener/internal/app/config"
	"github.com/AxMdv/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerConnector_CreateShortURL(t *testing.T) {
	type want struct {
		contentType string
		respBody    string
		statusCode  int
	}
	tests := []struct {
		name       string
		serC       *ShortenerHandlers
		requestURL string
		reqBody    string
		want       want
	}{
		{
			name:       "Positive test #1",
			serC:       &ShortenerHandlers{Repository: &storage.Repository{MapURL: map[string][]byte{}}},
			requestURL: "/",
			reqBody:    "https://yandex.ru",
			want: want{
				contentType: "text/plain",
				respBody:    "http://localhost:8080/aHR0cHM6Ly95YW5kZXgucnU",
				statusCode:  201,
			},
		},
	}
	config.ParseOptions()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			reqBody := bytes.NewReader([]byte(tt.reqBody))
			request := httptest.NewRequest(http.MethodPost, tt.requestURL, reqBody)
			w := httptest.NewRecorder()
			tt.serC.CreateShortURL(w, request)
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

func TestServerConnector_GetLongURL(t *testing.T) {
	type want struct {
		respHeaderLocation string
		statusCode         int
	}
	tests := []struct {
		name        string
		s           *ShortenerHandlers
		requestPath string
		want        want
	}{
		{
			name: "Positive test #1",
			s: &ShortenerHandlers{Repository: &storage.Repository{MapURL: map[string][]byte{
				"aHR0cHM6Ly95YW5kZXgucnU": []byte("https://yandex.ru")}},
			},
			requestPath: "/aHR0cHM6Ly95YW5kZXgucnU",
			want: want{
				respHeaderLocation: "https://yandex.ru",
				statusCode:         307,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get("/{shortenedURL}", tt.s.GetLongURL)
			ts := httptest.NewServer(r)
			defer ts.Close()

			request, err := http.NewRequest(http.MethodGet, ts.URL+(tt.requestPath), nil)
			require.NoError(t, err)
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			resp, err := client.Do(request)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.respHeaderLocation, resp.Header.Get("Location"))
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

		})
	}
}

func TestShortenerHandlers_CreateShortURLJson(t *testing.T) {

	type want struct {
		contentType string
		respBody    string
		statusCode  int
	}
	tests := []struct {
		name           string
		s              *ShortenerHandlers
		requestURL     string
		reqBody        string
		reqContentType string
		want           want
	}{
		{
			name:           "Positive test #1",
			s:              &ShortenerHandlers{Repository: &storage.Repository{MapURL: map[string][]byte{}}},
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
			tt.s.CreateShortURLJson(w, request)
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
