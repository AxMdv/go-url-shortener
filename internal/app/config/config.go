package config

import (
	"flag"
	"os"
)

var Options struct {
	RunAddr            string
	ResponseResultAddr string
	FileStorage        string
}

func ParseOptions() {
	flag.StringVar(&Options.RunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&Options.ResponseResultAddr, "b", "http://localhost:8080", "resut basic response address (before shortened URL)")
	flag.StringVar(&Options.FileStorage, "f", "/tmp/short-url-db.json", "path to save shortened URLs")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		Options.RunAddr = envRunAddr
	}
	if envResultAddr := os.Getenv("BASE_URL"); envResultAddr != "" {
		Options.ResponseResultAddr = envResultAddr
	}
	if envFileStorage := os.Getenv("FILE_STORAGE_PATH"); envFileStorage != "" {
		Options.FileStorage = envFileStorage
	}
}
