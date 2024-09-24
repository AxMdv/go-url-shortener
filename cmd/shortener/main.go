package main

import (
	"log"

	"github.com/AxMdv/go-url-shortener/internal/app"
)

func main() {

	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}

}
