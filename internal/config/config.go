package config

import (
	"flag"
	"os"
)

type Options struct {
	RunAddr            string
	ResponseResultAddr string
	FileStorage        string
	DataBaseDSN        string
}

func ParseOptions() *Options {
	Options := Options{}
	flag.StringVar(&Options.RunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&Options.ResponseResultAddr, "b", "http://localhost:8080", "resut basic response address (before shortened URL)")
	flag.StringVar(&Options.FileStorage, "f", "/tmp/short-url-db.json", "path to save shortened URLs")
	flag.StringVar(&Options.DataBaseDSN, "d", "", "dsn for acees to DB")
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
	if envDataBaseDSN := os.Getenv("DATABASE_DSN"); envDataBaseDSN != "" {
		Options.DataBaseDSN = envDataBaseDSN
	}
	return &Options
}

//dsn := "user=postgres password=adm dbname=postgres host=localhost port=5432 sslmode=disable"
