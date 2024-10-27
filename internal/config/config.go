// Package config describes options to run the app.
package config

import (
	"flag"
	"os"
)

// Options is parameters of running applications.
type Options struct {
	// RunAddr is the address and port to run server.
	RunAddr string
	// Resut basic response address (before shortened URL).
	ResponseResultAddr string
	// Path to save shortened URLs.
	FileStorage string
	// DSN for acees to DB.
	DataBaseDSN string
	// Enable HTTPS
	EnableHTTPS bool
}

// ParseOptions parses cmd flags and os environment variables.
func ParseOptions() *Options {
	options := Options{}
	flag.StringVar(&options.RunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&options.ResponseResultAddr, "b", "http://localhost:8080", "resut basic response address (before shortened URL)")
	flag.StringVar(&options.FileStorage, "f", "/tmp/short-url-db.json", "path to save shortened URLs")
	flag.StringVar(&options.DataBaseDSN, "d", "", "dsn for acees to DB")
	flag.BoolVar(&options.EnableHTTPS, "s", false, "enable https")
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
	return &options
}

// dsn := "user=postgres password=adm dbname=postgres host=localhost port=5432 sslmode=disable"
