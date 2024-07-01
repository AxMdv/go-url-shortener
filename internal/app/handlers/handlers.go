package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AxMdv/go-url-shortener/internal/app/config"
	"github.com/AxMdv/go-url-shortener/internal/app/model"
	"github.com/AxMdv/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

type ShortenerHandlers struct {
	Repository *storage.Repository
}

func (s *ShortenerHandlers) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	longURL, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	shortenedURL := base64.RawStdEncoding.EncodeToString(longURL)
	s.Repository.AddURL(longURL, shortenedURL)
	res := fmt.Sprintf("%s/%s", config.Options.ResponseResultAddr, shortenedURL)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func (s *ShortenerHandlers) GetLongURL(w http.ResponseWriter, r *http.Request) {
	shortenedURL := chi.URLParam(r, "shortenedURL")
	longURL, found := s.Repository.FindShortenedURL(shortenedURL)
	if !found {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.Header().Add("Location", string(longURL))
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func (s *ShortenerHandlers) CreateShortURLJson(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
		return
	}
	var request model.Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shortenedURL := base64.RawStdEncoding.EncodeToString([]byte(request.URL))
	s.Repository.AddURL([]byte(request.URL), shortenedURL)
	response := model.Response{Result: shortenedURL}
	resp, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
