package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"errors"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/AxMdv/go-url-shortener/internal/service"
	"github.com/AxMdv/go-url-shortener/pkg/auth"

	"github.com/AxMdv/go-url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

type ShortenerHandlers struct {
	shortenerService service.ShortenerService
	Config           config.Options
}

func NewShortenerHandlers(shortenerService service.ShortenerService, config *config.Options) *ShortenerHandlers {
	return &ShortenerHandlers{shortenerService: shortenerService, Config: *config}
}

func (s *ShortenerHandlers) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	longURL, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	shortenedURL := s.shortenerService.ShortenLongURL(longURL)

	formedURL := &storage.FormedURL{
		UUID:         auth.GetUUIDFromContext(r.Context()),
		ShortenedURL: shortenedURL,
		LongURL:      string(longURL),
	}
	err = s.shortenerService.CreateShortURL(formedURL)
	if err != nil {
		var duplicateErr *storage.AddURLError
		if errors.As(err, &duplicateErr) {
			res := fmt.Sprintf("%s/%s", s.Config.ResponseResultAddr, duplicateErr.DuplicateValue)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(res))
			return
		}
		log.Panic("Cant save urls to storage ", err)
		return
	}

	res := fmt.Sprintf("%s/%s", s.Config.ResponseResultAddr, shortenedURL)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func (s *ShortenerHandlers) GetLongURL(w http.ResponseWriter, r *http.Request) {
	shortenedURL := chi.URLParam(r, "shortenedURL")
	deleted, err := s.shortenerService.GetFlagByShortURL(shortenedURL)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if deleted {
		w.WriteHeader(http.StatusGone)
		return
	}
	longURL, err := s.shortenerService.GetLongURL(shortenedURL)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if longURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.Header().Add("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}

func (s *ShortenerHandlers) CreateShortURLJson(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
		return
	}
	var request Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortenedURL := s.shortenerService.ShortenLongURL([]byte(request.URL))
	formedURL := &storage.FormedURL{
		UUID:         auth.GetUUIDFromContext(r.Context()),
		ShortenedURL: shortenedURL,
		LongURL:      request.URL,
	}

	err := s.shortenerService.CreateShortURL(formedURL)
	if err != nil {
		var duplicateErr *storage.AddURLError
		if errors.As(err, &duplicateErr) {

			res := fmt.Sprintf("%s/%s", s.Config.ResponseResultAddr, duplicateErr.DuplicateValue)
			response := Response{Result: res}
			resp, err := json.Marshal(response)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(resp))
			return

		}
		log.Panic("Cant save urls to storage", err)
		return
	}

	res := fmt.Sprintf("%s/%s", s.Config.ResponseResultAddr, shortenedURL)
	response := Response{Result: res}
	resp, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (s *ShortenerHandlers) CheckDatabaseConnection(w http.ResponseWriter, r *http.Request) {
	if s.Config.DataBaseDSN == "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := s.shortenerService.PingDatabase()
	// err := storage.PingDatabase(s.Config)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *ShortenerHandlers) CreateShortURLBatch(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(bodyBytes) < 1 {
		w.WriteHeader(http.StatusNotAcceptable)
	}

	var requestBatch RequestBatch
	err = json.Unmarshal(bodyBytes, &requestBatch.BatchList)
	if err != nil {
		log.Panic(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	uuid := auth.GetUUIDFromContext(r.Context())
	formedURL := requestBatch.ToFormed(uuid)

	err = s.shortenerService.CreateShortURLBatch(formedURL)
	if err != nil {
		log.Panic("can`t add url batch to storage", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respData := make([]BatchShortened, len(formedURL))
	for i, v := range formedURL {
		respData[i].CorrelationID = v.CorrelationID
		respData[i].ShortenedURL = fmt.Sprintf("%s/%s", s.Config.ResponseResultAddr, v.ShortenedURL)
	}
	resp, err := json.Marshal(respData)
	if err != nil {
		log.Panic("can`t marshal response batch", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (s *ShortenerHandlers) GetAllURLByID(w http.ResponseWriter, r *http.Request) {

	uuid := auth.GetUUIDFromContext(r.Context())
	formedURL, err := s.shortenerService.GetAllURLByID(uuid)
	if err != nil {
		var noContentError *storage.NoContentError
		if errors.As(err, &noContentError) {
			w.WriteHeader(http.StatusNoContent)
		}
		w.WriteHeader(http.StatusInternalServerError)

	}
	for i, v := range formedURL {
		formedURL[i].ShortenedURL = fmt.Sprintf("%s/%s", s.Config.ResponseResultAddr, v.ShortenedURL)
	}
	resp, err := json.Marshal(formedURL)
	if err != nil {
		log.Panic("can`t marshal user urls", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

func (s *ShortenerHandlers) DeleteURLBatch(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
		return
	}

	uuid := auth.GetUUIDFromContext(r.Context())

	var deleteBatch storage.DeleteBatch

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(bodyBytes, &deleteBatch.ShortenedURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	deleteBatch.UUID = uuid

	s.shortenerService.DeleteURLBatch(deleteBatch)
	w.WriteHeader(http.StatusAccepted)
}
