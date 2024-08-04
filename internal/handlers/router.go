package handlers

import (
	mw "github.com/AxMdv/go-url-shortener/pkg/middleware"
	"github.com/go-chi/chi/v5"
)

func NewShortenerRouter(s *ShortenerHandlers) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", mw.WithLogging(mw.SignUpMiddleware(mw.GzipMiddleware(s.CreateShortURL))))
	r.Get("/{shortenedURL}", mw.WithLogging(s.GetLongURL))
	r.Post("/api/shorten", mw.WithLogging(mw.SignUpMiddleware(mw.GzipMiddleware(s.CreateShortURLJson))))
	r.Get("/ping", mw.WithLogging(s.CheckDatabaseConnection))
	r.Post("/api/shorten/batch", mw.WithLogging(mw.SignUpMiddleware(s.CreateShortURLBatch)))
	r.Get("/api/user/urls", mw.WithLogging(mw.ValidateUserMiddleware(mw.GzipMiddleware((s.GetAllURLByID)))))
	r.Delete("/api/user/urls", mw.WithLogging(mw.ValidateUserMiddleware(mw.GzipMiddleware(s.DeleteURLBatch))))
	return r
}
