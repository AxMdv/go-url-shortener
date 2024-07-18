package main

import (
	"log"
	"net/http"

	"github.com/AxMdv/go-url-shortener/pkg/logger"

	"github.com/AxMdv/go-url-shortener/internal/app/config"
	"github.com/AxMdv/go-url-shortener/internal/app/handlers"
	"github.com/AxMdv/go-url-shortener/pkg/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := config.ParseOptions()

	s, err := handlers.NewShortenerHandlers(cfg)
	if err != nil {
		log.Panic(err)
	}
	err = logger.InitLogger()
	if err != nil {
		log.Panic("Failed to init logger ", err)
	}
	r := chi.NewRouter()
	r.Post("/", middleware.WithLogging(middleware.GzipMiddleware(s.CreateShortURL)))
	r.Get("/{shortenedURL}", middleware.WithLogging(s.GetLongURL))
	r.Post("/api/shorten", middleware.WithLogging(middleware.GzipMiddleware(s.CreateShortURLJson)))
	r.Get("/ping", middleware.WithLogging(s.CheckDatabaseConnection))

	log.Fatal(http.ListenAndServe(cfg.RunAddr, r))
}
