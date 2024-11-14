package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/AxMdv/go-url-shortener/internal/service"
	"github.com/AxMdv/go-url-shortener/internal/storage"
)

func TestCreateShortURL(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        "",
	}
	repository, err := storage.NewRepository(config)
	require.NoError(t, err)
	urlService := service.NewShortenerService(repository)
	shortenerHandlers := NewShortenerHandlers(urlService, config)

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

func TestShortenerHandlersGetLongURL(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        "",
	}
	repository, err := storage.NewRepository(config)
	require.NoError(t, err)
	urlService := service.NewShortenerService(repository)
	shortenerHandlers := NewShortenerHandlers(urlService, config)

	router := NewShortenerRouter(shortenerHandlers)
	require.NoError(t, err)

	formedURL := storage.FormedURL{
		LongURL:      "https://yandex.ru",
		ShortenedURL: "aHR0cHM6Ly95YW5kZXgucnU",
		UUID:         "asd",
	}
	err = shortenerHandlers.shortenerService.CreateShortURL(&formedURL)
	require.NoError(t, err)
	type want struct {
		longURL    string
		statusCode int
	}
	tests := []struct {
		name       string
		requestURL string
		want       want
	}{
		{
			name:       "Positive test #1",
			requestURL: "http://localhost:8080/aHR0cHM6Ly95YW5kZXgucnU",

			want: want{
				longURL:    "https://yandex.ru",
				statusCode: 307,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(router)
			defer server.Close()

			request, err := http.NewRequest(http.MethodGet, server.URL+`/`+strings.TrimPrefix(tt.requestURL, "http://localhost:8080/"), nil)
			require.NoError(t, err)
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			resp, err := client.Do(request)
			require.NoError(t, err)
			defer resp.Body.Close()
			longURL := resp.Header.Get("Location")
			fmt.Println(longURL)

			assert.Equal(t, tt.want.longURL, longURL)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
		})

	}
}

func TestShortenerHandlersCreateShortURLJson(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        "",
	}
	repository, err := storage.NewRepository(config)
	require.NoError(t, err)
	urlService := service.NewShortenerService(repository)
	shortenerHandlers := NewShortenerHandlers(urlService, config)

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

func TestShortenerHandlersCheckDatabaseConnection(t *testing.T) {

	type want struct {
		statusCode int
	}
	tests := []struct {
		name        string
		requestURL  string
		databaseDSN string
		want        want
	}{
		{
			name:        "Negative test #1",
			requestURL:  "http://localhost:8080/ping",
			databaseDSN: "",
			want: want{
				statusCode: 405,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &config.Options{
				RunAddr:            ":8080",
				ResponseResultAddr: "http://localhost:8080",
				FileStorage:        "",
				DataBaseDSN:        "",
			}
			repository, err := storage.NewRepository(config)
			require.NoError(t, err)
			urlService := service.NewShortenerService(repository)
			shortenerHandlers := NewShortenerHandlers(urlService, config)

			request := httptest.NewRequest(http.MethodGet, tt.requestURL, nil)
			w := httptest.NewRecorder()
			shortenerHandlers.CheckDatabaseConnection(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestShortenerHandlersCreateShortURLBatch(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        "",
	}
	repository, err := storage.NewRepository(config)
	require.NoError(t, err)
	urlService := service.NewShortenerService(repository)
	shortenerHandlers := NewShortenerHandlers(urlService, config)

	router := NewShortenerRouter(shortenerHandlers)
	require.NoError(t, err)

	// formedURL := storage.FormedURL{
	// 	LongURL:       "https://yandex.ru",
	// 	ShortenedURL:  "aHR0cHM6Ly95YW5kZXgucnU",
	// 	CorrelationID: "123",
	// }
	// err = shortenerHandlers.shortenerService.CreateShortURL(&formedURL)
	// require.NoError(t, err)
	// formedURL = storage.FormedURL{
	// 	LongURL:       "https://vk.com",
	// 	ShortenedURL:  "aHR0cHM6Ly92ay5jb20",
	// 	CorrelationID: "321",
	// }
	// err = shortenerHandlers.shortenerService.CreateShortURL(&formedURL)
	// require.NoError(t, err)

	type want struct {
		responseBatch []BatchShortened
		statusCode    int
	}
	tests := []struct {
		name         string
		requestURL   string
		requestBatch []BatchOriginal
		want         want
	}{
		{
			name:       "Positive test #1",
			requestURL: "http://localhost:8080/api/shorten/batch",
			requestBatch: []BatchOriginal{
				{
					CorrelationID: "321",
					OriginalURL:   "https://vk.com",
				},
				{
					CorrelationID: "123",
					OriginalURL:   "https://yandex.ru",
				},
			},

			want: want{
				responseBatch: []BatchShortened{
					{
						CorrelationID: "321",
						ShortenedURL:  "http://localhost:8080/aHR0cHM6Ly92ay5jb20",
					},
					{
						CorrelationID: "123",
						ShortenedURL:  "http://localhost:8080/aHR0cHM6Ly95YW5kZXgucnU",
					},
				},
				statusCode: 201,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := httptest.NewServer(router)
			defer server.Close()
			bodyBytes, err := json.Marshal(tt.requestBatch)
			require.NoError(t, err)
			body := bytes.NewBuffer(bodyBytes)
			request, err := http.NewRequest(http.MethodPost, server.URL+`/`+strings.TrimPrefix(tt.requestURL, "http://localhost:8080/"), body)
			require.NoError(t, err)
			client := &http.Client{}
			resp, err := client.Do(request)
			require.NoError(t, err)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()
			responseBatch := make([]BatchShortened, len(tt.requestBatch))
			err = json.Unmarshal(respBody, &responseBatch)
			require.NoError(t, err)

			assert.Equal(t, tt.want.responseBatch, responseBatch)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
		})
	}
}

func TestShortenerHandlersGetAllURLByID(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        "",
	}
	repository, err := storage.NewRepository(config)
	require.NoError(t, err)
	urlService := service.NewShortenerService(repository)
	shortenerHandlers := NewShortenerHandlers(urlService, config)

	router := NewShortenerRouter(shortenerHandlers)
	require.NoError(t, err)

	type want struct {
		responseBatch []storage.FormedURL
		statusCode    int
	}
	tests := []struct {
		name       string
		requestURL string
		addURLs    []string
		want       want
	}{
		{
			name:       "Positive test #1",
			requestURL: "http://localhost:8080/api/user/urls",
			addURLs:    []string{"https://yandex.ru", "https://vk.com"},
			want: want{
				responseBatch: []storage.FormedURL{
					{
						LongURL:      "https://vk.com",
						ShortenedURL: "http://localhost:8080/aHR0cHM6Ly92ay5jb20",
					},
					{
						LongURL:      "https://yandex.ru",
						ShortenedURL: "http://localhost:8080/aHR0cHM6Ly95YW5kZXgucnU",
					},
				},
				statusCode: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := httptest.NewServer(router)
			defer server.Close()
			// first add some urls by user with cookie user_id
			cookie := &http.Cookie{Name: "user_id", Value: "30316566613236612d656333362d366533362d383031362d303031353564623832353563f8c1885334e5c310884a9af35b2f189228001d67f322295417f17ff534a88b99"}
			for _, addURL := range tt.addURLs {
				body := bytes.NewBuffer([]byte(addURL))
				req, err := http.NewRequest(http.MethodPost, server.URL, body)
				require.NoError(t, err)
				req.AddCookie(cookie)
				req.Header.Set("Content-Type", "text/html")
				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				resp.Body.Close()
				fmt.Printf("%+v", req)
				require.Equal(t, 201, resp.StatusCode)
			}
			request, err := http.NewRequest(http.MethodGet, server.URL+`/`+strings.TrimPrefix(tt.requestURL, "http://localhost:8080/"), nil)
			request.AddCookie(cookie)
			require.NoError(t, err)
			client := &http.Client{}
			resp, err := client.Do(request)
			require.NoError(t, err)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()
			responseBatch := make([]storage.FormedURL, len(tt.want.responseBatch))
			err = json.Unmarshal(respBody, &responseBatch)
			require.NoError(t, err)

			assert.ObjectsAreEqualValues(tt.want.responseBatch, responseBatch)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
		})
	}
}

func TestShortenerHandlersDeleteURLBatch(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        "",
	}
	repository, err := storage.NewRepository(config)
	require.NoError(t, err)
	urlService := service.NewShortenerService(repository)
	shortenerHandlers := NewShortenerHandlers(urlService, config)

	router := NewShortenerRouter(shortenerHandlers)
	require.NoError(t, err)

	type want struct {
		statusCode int
	}
	tests := []struct {
		name        string
		contentType string
		requestURL  string
		deleteURLs  []string
		want        want
	}{
		{
			name:        "Positive test #1",
			contentType: "application/json",
			requestURL:  "http://localhost:8080/api/user/urls",
			deleteURLs:  []string{"http://localhost:8080/aHR0cHM6Ly92ay5jb20", "http://localhost:8080/aHR0cHM6Ly95YW5kZXgucnU"},
			want: want{
				statusCode: 202,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := httptest.NewServer(router)
			defer server.Close()
			cookie := &http.Cookie{Name: "user_id", Value: "30316566613236612d656333362d366533362d383031362d303031353564623832353563f8c1885334e5c310884a9af35b2f189228001d67f322295417f17ff534a88b99"}

			rBody, err := json.Marshal(tt.deleteURLs)
			require.NoError(t, err)
			body := bytes.NewBuffer(rBody)

			request, err := http.NewRequest(http.MethodDelete, server.URL+`/`+strings.TrimPrefix(tt.requestURL, "http://localhost:8080/"), body)
			request.AddCookie(cookie)
			request.Header.Set("Content-Type", tt.contentType)
			require.NoError(t, err)
			client := &http.Client{}
			resp, err := client.Do(request)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
		})
	}
}
