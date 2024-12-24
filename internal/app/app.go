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
	"time"

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
	fmt.Printf("%+v\n", a.configOptions)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	go func() {
		fmt.Println("waiting for ctrl+c")
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := a.server.Shutdown(shutdownCtx); err != nil {
			// ошибки закрытия Listener
			log.Printf("error in HTTP server Shutdown: %v\n", err)

		} else {
			log.Println("successfully stopped http server")
		}
	}()
	if err := a.runHTTPServer(); err != http.ErrServerClosed {
		log.Panicln(err)
	}
	err := a.gracefullShutdown()
	if err != nil {
		log.Print(err)
	}
	log.Println("shutting down...")
	os.Exit(0)
	return nil

}

func (a *App) runHTTPServer() error {

	if a.configOptions.EnableHTTPS {
		log.Printf("HTTPS server is running on %s", a.configOptions.RunAddr)
		return a.server.ListenAndServeTLS("./certs/certbundle.pem", "./certs/server.key")
	}
	log.Printf("HTTP server is running on %s", a.configOptions.RunAddr)
	return a.server.ListenAndServe()
}

func (a *App) gracefullShutdown() error {
	// close repo if it has method close()
	fmt.Println("trying to close repository..")
	_, ok := a.urlRepository.(Closer)
	var err error
	if ok {
		err = a.urlRepository.(Closer).Close()
		if err != nil {
			log.Println("error in closing repo", err)
			return err
		}
		log.Println("success in closing repo")
	} else {
		log.Println("current repo doesn`t have method Close()")
	}

	return err
}
