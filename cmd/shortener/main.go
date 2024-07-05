package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AxMdv/go-url-shortener/internal/app/compress"
	"github.com/AxMdv/go-url-shortener/internal/app/config"
	"github.com/AxMdv/go-url-shortener/internal/app/handlers"
	"github.com/AxMdv/go-url-shortener/internal/app/logger"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.ParseOptions()
	s, err := handlers.InitShortenerHandlers(config.Options.FileStorage)
	if err != nil {
		fmt.Println("Failed to init ShortenerHandlers", err)
	}
	err = logger.InitLogger()
	if err != nil {
		fmt.Println("Failed to init logger")
	}
	r := chi.NewRouter()
	r.Post("/", logger.WithLogging(compress.GzipMiddleware(s.CreateShortURL)))
	r.Get("/{shortenedURL}", logger.WithLogging(compress.GzipMiddleware(s.GetLongURL)))
	r.Post("/api/shorten", logger.WithLogging(compress.GzipMiddleware(s.CreateShortURLJson)))

	log.Fatal(http.ListenAndServe(config.Options.RunAddr, r))
}
