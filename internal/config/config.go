// Package config describes options to run the app.
package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

// Options is parameters of running applications.
type Options struct {
	ConfigPath string `json:"-"`
	// RunAddr is the address and port to run server.
	RunAddr string `json:"server_address"`
	// Resut basic response address (before shortened URL).
	ResponseResultAddr string `json:"base_url"`
	// Path to save shortened URLs.
	FileStorage string `json:"file_storage_path"`
	// DSN for acees to DB.
	DataBaseDSN string `json:"database_dsn"`
	// Enable HTTPS
	EnableHTTPS bool `json:"enable_https"`
}

// ParseOptions parses cmd flags and os environment variables.
func ParseOptions() *Options {
	options := Options{}

	flag.StringVar(&options.RunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&options.ResponseResultAddr, "b", "http://localhost:8080", "resut basic response address (before shortened URL)")
	flag.StringVar(&options.FileStorage, "f", "/tmp/short-url-db.json", "path to save shortened URLs")
	flag.StringVar(&options.DataBaseDSN, "d", "", "dsn for acees to DB")
	flag.BoolVar(&options.EnableHTTPS, "s", false, "enable https")
	flag.StringVar(&options.ConfigPath, "c", "", "path to config file")
	flag.StringVar(&options.ConfigPath, "config", "", "path to config file")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		options.RunAddr = envRunAddr
	}
	if envResultAddr := os.Getenv("BASE_URL"); envResultAddr != "" {
		options.ResponseResultAddr = envResultAddr
	}
	if envFileStorage := os.Getenv("FILE_STORAGE_PATH"); envFileStorage != "" {
		options.FileStorage = envFileStorage
	}
	if envDataBaseDSN := os.Getenv("DATABASE_DSN"); envDataBaseDSN != "" {
		options.DataBaseDSN = envDataBaseDSN
	}
	if envDataBaseDSN := os.Getenv("ENABLE_HTTPS"); envDataBaseDSN != "" {
		options.DataBaseDSN = envDataBaseDSN
	}
	if envConfigPath := os.Getenv("CONFIG"); envConfigPath != "" {
		options.ConfigPath = envConfigPath
	}
	confOpts := &Options{}
	if options.ConfigPath != "" {
		confFile, err := os.ReadFile(options.ConfigPath)
		if err != nil {
			log.Panic(err)
		}
		err = json.Unmarshal(confFile, confOpts)
		if err != nil {
			log.Panic(err)
		}
		if options.RunAddr == "" {
			options.RunAddr = confOpts.RunAddr
		}
		if options.ResponseResultAddr == "" {
			options.ResponseResultAddr = confOpts.ResponseResultAddr
		}
		if options.FileStorage == "" {
			options.FileStorage = confOpts.FileStorage
		}
		if options.DataBaseDSN == "" {
			options.DataBaseDSN = confOpts.DataBaseDSN
		}
		if !options.EnableHTTPS {
			options.EnableHTTPS = confOpts.EnableHTTPS
		}
	}

	return &options
}

// dsn := "user=postgres password=adm dbname=postgres host=localhost port=5432 sslmode=disable"

// set FILE_STORAGE_PATH=''
// shortenertestbeta -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener
// shortenertestbeta -test.v -test.run=^TestIteration12$ -binary-path=cmd/shortener/shortener -database-dsn='user=postgres password=adm dbname=postgres host=localhost port=5432 sslmode=disable'
