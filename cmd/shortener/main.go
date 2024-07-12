package main

import (
	"log"
	"net/http"

	"github.com/AxMdv/go-url-shortener/internal/app/config"
	"github.com/AxMdv/go-url-shortener/internal/app/handlers"
	"github.com/AxMdv/go-url-shortener/pkg/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.ParseOptions()
	s, err := handlers.NewShortenerHandlers(config.Options.FileStorage)
	if err != nil {
		log.Panic("Failed to init ShortenerHandlers", err)
	}
	err = middleware.InitLogger()
	if err != nil {
		log.Panic("Failed to init logger", err)
	}
	r := chi.NewRouter()
	r.Post("/", middleware.WithLogging(middleware.GzipMiddleware(s.CreateShortURL)))
	r.Get("/{shortenedURL}", middleware.WithLogging(middleware.GzipMiddleware(s.GetLongURL)))
	r.Post("/api/shorten", middleware.WithLogging(middleware.GzipMiddleware(s.CreateShortURLJson)))

	log.Fatal(http.ListenAndServe(config.Options.RunAddr, r))
}
