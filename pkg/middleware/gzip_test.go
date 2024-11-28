package middleware_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/AxMdv/go-url-shortener/internal/handlers"
	"github.com/AxMdv/go-url-shortener/internal/service"
	"github.com/AxMdv/go-url-shortener/internal/storage"
	"github.com/AxMdv/go-url-shortener/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestGzipCompression(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        "",
	}
	repository, err := storage.NewRepository(config)
	require.NoError(t, err)
	urlService := service.NewShortenerService(repository)
	shortenerHandlers := handlers.NewShortenerHandlers(urlService, config)
	mux := chi.NewMux()
	mux.HandleFunc("/", middleware.GzipMiddleware(shortenerHandlers.CreateShortURL))

	srv := httptest.NewServer(mux)
	defer srv.Close()

	requestBody1 := `https://practicum.yandex.ru`

	requestBody2 := `https://yandex.ru`
	// ожидаемое содержимое тела ответа при успешном запросе
	successBody1 := `http://localhost:8080/aHR0cHM6Ly9wcmFjdGljdW0ueWFuZGV4LnJ1`
	successBody2 := `http://localhost:8080/aHR0cHM6Ly95YW5kZXgucnU`
	t.Run("send_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody1))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Accept-Encoding", "")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, successBody1, string(b))
	})

	t.Run("accept_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString(requestBody2)
		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Accept-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		require.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))
		defer resp.Body.Close()

		zr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)

		b, err := io.ReadAll(zr)
		require.NoError(t, err)

		require.Equal(t, successBody2, string(b))
	})
}
