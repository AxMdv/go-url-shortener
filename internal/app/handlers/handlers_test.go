package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AxMdv/go-url-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerConnector_HandleMethod(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name          string
		serC          *ServerConnector
		requestURL    string
		requestMethod string
		want          want
	}{
		{
			name:          "Negative test #1",
			serC:          &ServerConnector{StC: &storage.StorageConnector{MapURL: map[string][]byte{}}},
			requestURL:    "/",
			requestMethod: "PUT",
			want: want{
				statusCode: 405,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.requestMethod, tt.requestURL, nil)
			w := httptest.NewRecorder()
			tt.serC.HandleMethod(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestServerConnector_HandlePostMain(t *testing.T) {
	type want struct {
		contentType string
		respBody    string
		statusCode  int
	}
	tests := []struct {
		name       string
		serC       *ServerConnector
		requestURL string
		reqBody    string
		want       want
	}{
		{
			name:       "Positive test #1",
			serC:       &ServerConnector{StC: &storage.StorageConnector{MapURL: map[string][]byte{}}},
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
			tt.serC.HandlePostMain(w, request)
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

func TestServerConnector_HandleGetMain(t *testing.T) {
	// type args struct {
	// 	w http.ResponseWriter
	// 	r *http.Request
	// }
	type want struct {
		respHeaderLocation string
		statusCode         int
	}
	tests := []struct {
		name       string
		serC       *ServerConnector
		requestURL string
		want       want
	}{
		{
			name: "Positive test #1",
			serC: &ServerConnector{StC: &storage.StorageConnector{MapURL: map[string][]byte{
				"aHR0cHM6Ly95YW5kZXgucnU": []byte("https://yandex.ru")}},
			},
			requestURL: "http://localhost:8080/aHR0cHM6Ly95YW5kZXgucnU",
			want: want{
				respHeaderLocation: "https://yandex.ru",
				statusCode:         307,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.requestURL, nil)
			w := httptest.NewRecorder()
			tt.serC.HandleGetMain(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.respHeaderLocation, result.Header.Get("Location"))
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func Test_shortenURL(t *testing.T) {
	type args struct {
		longURL []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shortenURL(tt.args.longURL); got != tt.want {
				t.Errorf("shortenURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
