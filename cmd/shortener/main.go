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
	serC := handlers.ServerConnector{StC: &storage.StorageConnector{MapURL: make(map[string][]byte)}}

	r := chi.NewRouter()
	r.Post("/", serC.HandlePostMain)
	r.Get("/{shortenedURL}", serC.HandleGetMain)

	log.Fatal(http.ListenAndServe(config.Options.RunAddr, r))
}
