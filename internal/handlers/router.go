package handlers

import (
	"net/http/pprof"

	"github.com/go-chi/chi/v5"

	mw "github.com/AxMdv/go-url-shortener/pkg/middleware"
)

// NewShortenerRouter returns mux, which handles all application api methods and profiler methods.
func NewShortenerRouter(s *ShortenerHandlers) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", mw.WithLogging(mw.SignUpMiddleware(mw.GzipMiddleware(s.CreateShortURL))))
	r.Get("/{shortenedURL}", mw.WithLogging(s.GetLongURL))
	r.Post("/api/shorten", mw.WithLogging(mw.SignUpMiddleware(mw.GzipMiddleware(s.CreateShortURLJson))))
	r.Get("/ping", mw.WithLogging(s.CheckDatabaseConnection))
	r.Post("/api/shorten/batch", mw.WithLogging(mw.SignUpMiddleware(mw.GzipMiddleware(s.CreateShortURLBatch))))
	r.Get("/api/user/urls", mw.WithLogging(mw.ValidateUserMiddleware(mw.GzipMiddleware((s.GetAllURLByID)))))
	r.Delete("/api/user/urls", mw.WithLogging(mw.ValidateUserMiddleware(mw.GzipMiddleware(s.DeleteURLBatch))))
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	return r
}
