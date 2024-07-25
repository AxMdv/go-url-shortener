package app

import (
	"log"
	"net/http"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/AxMdv/go-url-shortener/internal/handlers"
	"github.com/AxMdv/go-url-shortener/pkg/logger"
)

type App struct {
	serviceProvider *serviceProvider
}

func NewApp() (*App, error) {
	a := &App{}

	// err := a.initServiceProvider()
	// if err != nil {
	// 	return nil, err
	// }
	// err = a.initConfigOptions()
	// if err != nil {
	// 	return nil, err
	// }
	// a.initURLRepository()
	// if err != nil {
	// 	return nil, err
	// }
	err := a.initDeps()
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run() error {
	return a.runHTTPServer()
}
func (a *App) initDeps() error {
	inits := []func() error{
		a.initServiceProvider,
		a.initConfigOptions,
		a.initLogger,
		// a.initURLRepository,
		a.initShortenerHanders,
		a.initRouter,
	}

	for _, f := range inits {
		err := f()
		if err != nil {
			return err
		}
	}

	return nil
}
func (a *App) initServiceProvider() error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initConfigOptions() error {
	a.serviceProvider.configOptions = config.ParseOptions()
	return nil
}

func (a *App) initLogger() error {
	err := logger.InitLogger()
	return err
}

func (a *App) initShortenerHanders() error {
	a.serviceProvider.ShortenerHandlers()
	return nil
}

func (a *App) initRouter() error {
	a.serviceProvider.router = handlers.NewShortenerRouter(a.serviceProvider.shortenerHandlers)
	return nil
}

func (a *App) runHTTPServer() error {
	log.Printf("HTTP server is running on %s", a.serviceProvider.configOptions.RunAddr)
	return http.ListenAndServe(a.serviceProvider.configOptions.RunAddr, a.serviceProvider.router)
}

// func (a *App) initURLRepository() error {

// 	temp, err := storage.NewRepository(a.serviceProvider.configOptions) //!!!!!!!!!!
// 	a.serviceProvider.urlRepository = temp
// 	return err
// }
