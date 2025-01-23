package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseOptions(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		cfg, err := os.Create("config.json")
		require.NoError(t, err)
		defer os.Remove(cfg.Name())
		conf := `{
"server_address": "localhost:8080",
"base_url": "http://localhost",
"file_storage_path": "./tmp/short-url-db.json",
"database_dsn": "123",
"enable_https": false,
"trusted_subnet": "::1/128",
"grpc_server_address": ":8090"
} `
		cfg.Write([]byte(conf))
		opt := ParseOptions()
		assert.Equal(t, ":8080", opt.RunAddr)
		assert.Equal(t, "http://localhost:8080", opt.ResponseResultAddr)
		assert.Equal(t, "./tmp/short-url-db.json", opt.FileStorage)
		assert.Equal(t, "123", opt.DataBaseDSN)
	})

}

// func TestParseOptionsWithConfig(t *testing.T) {

// 	cfg, err := os.Create("config.json")
// 	require.NoError(t, err)
// 	defer os.Remove(`./` + cfg.Name())
// 	conf := `{
// "server_address": "localhost:8080",
// "base_url": "http://localhost",
// "file_storage_path": "./tmp/short-url-db.json",
// "database_dsn": "123",
// "enable_https": false,
// "trusted_subnet": "::1/128",
// "grpc_server_address": ":8090"
// } `

// 	opt := ParseOptions()
// 	assert.Equal(t, "", opt.DataBaseDSN)

// }
