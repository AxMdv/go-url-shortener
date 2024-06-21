package main

import (
	"net/http"

	"github.com/AxMdv/go-url-shortener/internal/app/handlers"
	"github.com/AxMdv/go-url-shortener/internal/app/storage"
)

func main() {
	serC := handlers.ServerConnector{StC: &storage.StorageConnector{MapURL: make(map[string][]byte)}}
	mux := http.NewServeMux()
	mux.HandleFunc("/", serC.HandleMethod)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
