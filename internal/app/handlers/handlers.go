package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/AxMdv/go-url-shortener/internal/app/config"
	"github.com/AxMdv/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

func (serC *ServerConnector) HandlePostMain(w http.ResponseWriter, r *http.Request) {
	longURL, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	shortenedURL := shortenURL(longURL)
	serC.StC.AddURL(longURL, shortenedURL)
	res := fmt.Sprintf("%s/%s", config.Options.ResponseResultAddr, shortenedURL)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func shortenURL(longURL []byte) string {
	res := base64.RawStdEncoding.EncodeToString(longURL)
	return res
}

func (serC *ServerConnector) HandleGetMain(w http.ResponseWriter, r *http.Request) {
	shortenedURL := chi.URLParam(r, "shortenedURL")
	longURL, found := serC.StC.FindShortenedURL(shortenedURL)
	if !found {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.Header().Add("Location", string(longURL))
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

type ServerConnector struct {
	StC *storage.StorageConnector
}
