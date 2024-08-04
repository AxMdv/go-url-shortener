package app

import (
	"log"
	"net/http"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/AxMdv/go-url-shortener/internal/handlers"
	"github.com/AxMdv/go-url-shortener/internal/service"
	"github.com/AxMdv/go-url-shortener/internal/storage"
	"github.com/AxMdv/go-url-shortener/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type App struct {
	urlService service.ShortenerService

	urlRepository storage.Repository

	configOptions *config.Options

	ShortenerHandlers *handlers.ShortenerHandlers

	router *chi.Mux
}

func NewApp() (*App, error) {

	config := config.ParseOptions()

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

func (a *App) Run() error {
	return a.runHTTPServer()
}

func (a *App) runHTTPServer() error {
	log.Printf("HTTP server is running on %s", a.configOptions.RunAddr)
	return http.ListenAndServe(a.configOptions.RunAddr, a.router)
}
