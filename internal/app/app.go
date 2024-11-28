// Package app is the implementation of shortener app.
package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/AxMdv/go-url-shortener/internal/handlers"
	"github.com/AxMdv/go-url-shortener/internal/service"
	"github.com/AxMdv/go-url-shortener/internal/storage"
	"github.com/AxMdv/go-url-shortener/pkg/logger"
)

// App is an application of url shortener
type App struct {
	urlService handlers.IShortenerService

	urlRepository service.IRepository

	configOptions *config.Options
	// ShortenerHandlers is api handlers
	ShortenerHandlers *handlers.ShortenerHandlers

	router *chi.Mux

	server *http.Server
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

	srv := &http.Server{
		Addr:    config.RunAddr,
		Handler: router,
	}

	a := &App{
		urlService:        urlService,
		urlRepository:     repository,
		configOptions:     config,
		ShortenerHandlers: shortenerHandlers,
		router:            router,
		server:            srv,
	}
	return a, nil
}

// Run is a main process of working application
func (a *App) Run() error {
	idleConnsClosed := make(chan struct{})

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go a.processInterrupt(sigint, idleConnsClosed)

	if err := a.runHTTPServer(); err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed
	err := a.gracefullShutdown()
	if err != nil {
		fmt.Print(err)
	}
	return err
}

func (a *App) runHTTPServer() error {

	if a.configOptions.EnableHTTPS {
		log.Printf("HTTPS server is running on %s", a.configOptions.RunAddr)
		return a.server.ListenAndServeTLS("./certs/certbundle.pem", "./certs/server.key")
	}
	log.Printf("HTTP server is running on %s", a.configOptions.RunAddr)
	return a.server.ListenAndServe()
}

func (a *App) processInterrupt(sigint chan os.Signal, idleConnsClosed chan struct{}) {
	<-sigint
	if err := a.server.Shutdown(context.Background()); err != nil {
		// ошибки закрытия Listener
		log.Printf("HTTP server Shutdown: %v", err)
	} else {
		log.Println("successfully stopped http server")
	}
	close(idleConnsClosed)
}

func (a *App) gracefullShutdown() error {
	// close repo if it has method close()
	log.Println("trying to close repository..")
	closerRepo, ok := a.urlRepository.(service.Closer)
	var err error
	if ok {
		err = closerRepo.Close()
		if err == nil {
			log.Println("success in closing repo")
		}
	}
	return err
}
