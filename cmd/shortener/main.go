package main

import (
	"log"

	"github.com/AxMdv/go-url-shortener/internal/app"
)

func main() {

	// cfg := config.ParseOptions()
	// s, err := handlers.NewShortenerHandlers(cfg)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// err = logger.InitLogger()
	// if err != nil {
	// 	log.Panic("Failed to init logger ", err)
	// }
	// router := handlers.NewShortenerRouter(s)
	// log.Fatal(http.ListenAndServe(cfg.RunAddr, router))
	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}
	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}

}
