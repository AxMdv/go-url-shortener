package handlers

import (
	"github.com/AxMdv/go-url-shortener/pkg/middleware"
	"github.com/go-chi/chi/v5"
)

func NewShortenerRouter(s *ShortenerHandlers) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", middleware.WithLogging(middleware.GzipMiddleware(s.CreateShortURL)))
	r.Get("/{shortenedURL}", middleware.WithLogging(s.GetLongURL))
	r.Post("/api/shorten", middleware.WithLogging(middleware.GzipMiddleware(s.CreateShortURLJson)))
	r.Get("/ping", middleware.WithLogging(s.CheckDatabaseConnection))
	r.Post("/api/shorten/batch", middleware.WithLogging(s.CreateShortURLBatch))
	return r
}
