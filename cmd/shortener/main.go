package main

import (
	"log"
	"net/http"

	"github.com/AxMdv/go-url-shortener/internal/app/config"
	"github.com/AxMdv/go-url-shortener/internal/app/handlers"
	"github.com/AxMdv/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.ParseOptions()
	s := handlers.ShortenerHandlers{R: &storage.Repository{MapURL: make(map[string][]byte)}}

	r := chi.NewRouter()
	r.Post("/", s.HandlePostMain)
	r.Get("/{shortenedURL}", s.HandleGetMain)

	log.Fatal(http.ListenAndServe(config.Options.RunAddr, r))
}
