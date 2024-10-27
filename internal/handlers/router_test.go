package handlers

import (
	"testing"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/AxMdv/go-url-shortener/internal/service"
	"github.com/AxMdv/go-url-shortener/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestNewShortenerRouter(t *testing.T) {
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
	_ = NewShortenerRouter(shortenerHandlers)
}
