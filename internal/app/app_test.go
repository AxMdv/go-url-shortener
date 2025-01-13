package app

import (
	"testing"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/stretchr/testify/require"
)

func TestNewApp(t *testing.T) {
	cfg := &config.Options{}
	_, err := NewApp(cfg)
	require.NoError(t, err)
}

// func TestAppRun(t *testing.T) {
// 	cfg := &config.Options{
// 		RunAddr: ":8080",
// 	}
// 	app, err := NewApp(cfg)
// 	require.NoError(t, err)

// 	err = app.Run()
// 	require.NoError(t, err)
// }
