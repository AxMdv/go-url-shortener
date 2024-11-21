// Package storage is designed to create requests to database.
package storage

import (
	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/AxMdv/go-url-shortener/internal/service"
)

// NewRepository returns repository, that implements service.IRepository interface.
func NewRepository(config *config.Options) (service.IRepository, error) {
	if config.DataBaseDSN != "" {
		return NewDBRepository(config)
	}
	if config.FileStorage != "" {
		return NewFileRepository(config)
	}
	return NewRAMRepository()
}

// .............................................................
