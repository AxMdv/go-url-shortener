// Package app is the implementation of shortener app.
package app

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/AxMdv/go-url-shortener/internal/handlers"
	"github.com/AxMdv/go-url-shortener/internal/service"
	"github.com/AxMdv/go-url-shortener/internal/storage"
	"github.com/AxMdv/go-url-shortener/pkg/logger"
)

// App is an application of url shortener
type App struct {
	urlService service.ShortenerService

	urlRepository storage.Repository

	configOptions *config.Options
	// ShortenerHandlers is api handlers
	ShortenerHandlers *handlers.ShortenerHandlers

	router *chi.Mux
}

// NewApp creates a new app of a URL shortener
func NewApp(config *config.Options) (*App, error) {

	err := logger.InitLogger()
	if err != nil {
		return nil, err
	}

	repository, err := storage.NewRepository(config)
	if err != nil {
		return nil, err
	}

	urlService := service.NewShortenerService(repository)

	shortenerHandlers := handlers.NewShortenerHandlers(urlService, config)

	router := handlers.NewShortenerRouter(shortenerHandlers)

	a := &App{
		urlService:        urlService,
		urlRepository:     repository,
		configOptions:     config,
		ShortenerHandlers: shortenerHandlers,
		router:            router,
	}
	return a, nil
}

// Run is a main process of working application
func (a *App) Run() error {
	return a.runHTTPServer()
}

func (a *App) runHTTPServer() error {
	if a.configOptions.EnableHTTPS {
		log.Printf("HTTPS server is running on %s", a.configOptions.RunAddr)
		return http.ListenAndServeTLS(a.configOptions.RunAddr, "./certs/certbundle.pem", "./certs/server.key", a.router)
	}
	log.Printf("HTTP server is running on %s", a.configOptions.RunAddr)
	return http.ListenAndServe(a.configOptions.RunAddr, a.router)
}
