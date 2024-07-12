package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AxMdv/go-url-shortener/internal/app/config"
	"github.com/AxMdv/go-url-shortener/internal/app/handlers"
	"github.com/AxMdv/go-url-shortener/internal/app/logger"
	"github.com/AxMdv/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.ParseOptions()
	s := handlers.ShortenerHandlers{Repository: &storage.Repository{MapURL: make(map[string][]byte)}}

	err := logger.InitLogger()
	if err != nil {
		fmt.Println("Failed to init logger")
	}
	r := chi.NewRouter()
	r.Post("/", logger.WithLogging(s.CreateShortURL))
	r.Get("/{shortenedURL}", logger.WithLogging(s.GetLongURL))
	r.Post("/api/shorten", logger.WithLogging(s.CreateShortURLJson))

	log.Fatal(http.ListenAndServe(config.Options.RunAddr, r))
}
