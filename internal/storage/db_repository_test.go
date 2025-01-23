package storage

import (
	"testing"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/stretchr/testify/require"
)

// SetDSNForTests sets dsn for postgres
func setDSNForTests() string {
	return "user=postgres password=adm dbname=postgres host=localhost port=5432 sslmode=disable"
}
func TestNewDBRepository(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        setDSNForTests(),
	}
	_, err := NewDBRepository(config)
	require.NoError(t, err)
}
