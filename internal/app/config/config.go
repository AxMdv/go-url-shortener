package config

import (
	"flag"
	"os"
)

var Options struct {
	RunAddr            string
	ResponseResultAddr string
}

func ParseOptions() {
	flag.StringVar(&Options.RunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&Options.ResponseResultAddr, "b", "http://localhost:8080", "resut basic response address (before shortened URL)")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		Options.RunAddr = envRunAddr
	}
	if envResultAddr := os.Getenv("BASE_URL"); envResultAddr != "" {
		Options.ResponseResultAddr = envResultAddr
	}
}
