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

type ShortenerHandlers struct {
	R *storage.Repository
}

func (s *ShortenerHandlers) HandlePostMain(w http.ResponseWriter, r *http.Request) {
	longURL, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	shortenedURL := shortenURL(longURL)
	s.R.AddURL(longURL, shortenedURL)
	res := fmt.Sprintf("%s/%s", config.Options.ResponseResultAddr, shortenedURL)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func shortenURL(longURL []byte) string {
	res := base64.RawStdEncoding.EncodeToString(longURL)
	return res
}

func (s *ShortenerHandlers) HandleGetMain(w http.ResponseWriter, r *http.Request) {
	shortenedURL := chi.URLParam(r, "shortenedURL")
	longURL, found := s.R.FindShortenedURL(shortenedURL)
	if !found {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.Header().Add("Location", string(longURL))
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
