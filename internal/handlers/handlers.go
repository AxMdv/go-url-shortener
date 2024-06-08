package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/AxMdv/go-url-shortener/internal/storage"
)

func (serC *ServerConnector) HandleMethod(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		serC.HandleGetMain(w, r)
	case http.MethodPost:
		serC.HandlePostMain(w, r)
	default:
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
}

func (serC *ServerConnector) HandlePostMain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	longURL, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer r.Body.Close()
	shortenedURL := shortenURL(longURL)
	serC.StC.AddURL(longURL, shortenedURL)
	res := fmt.Sprintf(`http://localhost:8080/%s`, shortenedURL)
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func shortenURL(longURL []byte) string {
	res := base64.RawStdEncoding.EncodeToString(longURL)
	return res
}

func (serC *ServerConnector) HandleGetMain(w http.ResponseWriter, r *http.Request) {
	shortenedURL := r.URL.Path[1:]
	longURL, found := serC.StC.FindShortenedURL(shortenedURL)
	if found {
		w.Header().Add("Location", string(longURL))
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

type ServerConnector struct {
	StC *storage.StorageConnector
}
