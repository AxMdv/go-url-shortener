package storage

import (
	"testing"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/stretchr/testify/require"
)

func TestNewDBRepository(t *testing.T) {
	config := &config.Options{
		DataBaseDSN: "user=postgres password=adm dbname=postgres host=localhost port=5432 sslmode=disable",
	}
	dr, err := NewDBRepository(config)
	require.NoError(t, err)
	defer dr.Close()
}
