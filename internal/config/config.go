// Package config describes options to run the app.
package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"runtime"
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
	// CIDR
	TrustedSubnet string `json:"trusted_subnet"`
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
	flag.StringVar(&options.TrustedSubnet, "t", "", "subnet to access to /api/internal/stats")
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
	if envTrustedSubnet := os.Getenv("TRUSTED_SUBNET"); envTrustedSubnet != "" {
		options.TrustedSubnet = envTrustedSubnet
	}
	confOpts := &Options{}
	if options.ConfigPath != "" {
		confFile, err := os.ReadFile(options.ConfigPath)
		if err != nil {
			log.Fatal(err)
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
	// adding dot in filepath on windows
	if options.FileStorage != "" && runtime.GOOS == "windows" {
		options.FileStorage = `.` + options.FileStorage
	}
	return &options
}

// dsn := "user=postgres password=adm dbname=postgres host=localhost port=5432 sslmode=disable"
// -database-dsn='postgresql://postgres:adm@127.0.0.1:5432/postgres?sslmode=disable'
// -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable'
// shortenertestbeta -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener
//-database-dsn='postgresql://postgres:adm@127.0.0.1:5432/postgres?sslmode=disable'

// shortenertest -test.v -test.run=^TestIteration11$ -binary-path=cmd/shortener/shortener -database-dsn='postgresql://postgres:adm@127.0.0.1:5432/postgres?sslmode=disable'
// shortenertest -test.v -test.run=^TestIteration11$ -binary-path=cmd/shortener/shortener -file-storage-path="./tmp/short-url-db.json" -server-port="8080" -source-path="." -database-dsn="user=postgres password=adm dbname=postgres host=localhost port=5432 sslmode=disable"
