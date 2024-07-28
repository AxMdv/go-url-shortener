package app

import (
	"log"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/AxMdv/go-url-shortener/internal/handlers"
	"github.com/AxMdv/go-url-shortener/internal/service"
	"github.com/AxMdv/go-url-shortener/internal/service/url"
	"github.com/AxMdv/go-url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

type serviceProvider struct {
	urlService service.ShortenerService

	urlRepository storage.Repository

	configOptions *config.Options

	shortenerHandlers *handlers.ShortenerHandlers

	router *chi.Mux
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) ShortenerHandlers() *handlers.ShortenerHandlers {
	if s.shortenerHandlers == nil {

		s.shortenerHandlers = handlers.NewShortenerHandlers(s.URLService())
		s.shortenerHandlers.Config = *s.configOptions
	}
	return s.shortenerHandlers
}

func (s *serviceProvider) URLService() service.ShortenerService {
	if s.urlService == nil {
		s.urlService = url.NewShortenerService(s.URLRepository())
	}
	return s.urlService
}

func (s *serviceProvider) URLRepository() storage.Repository {
	if s.urlRepository == nil {
		temp, err := storage.NewRepository(s.configOptions)
		if err != nil {
			log.Panicf("Failed to create repository, %v", err)
			return nil
		}
		s.urlRepository = temp
	}
	return s.urlRepository
}
