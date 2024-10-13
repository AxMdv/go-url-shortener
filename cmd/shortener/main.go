package main

import (
	"fmt"
	"log"

	"github.com/AxMdv/go-url-shortener/internal/app"
	"github.com/AxMdv/go-url-shortener/internal/config"
)

// go run -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X main.buildCommit=hello world" main.go
// go run -ldflags "-X main.buildVersion=v1.0.1 -X main.buildCommit=hello-world" main.go -f=""

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func printBuildInfo() {
	fmt.Printf("Build version:%s\n", formatValue(buildVersion))
	fmt.Printf("Build date:%s\n", formatValue(buildDate))
	fmt.Printf("Build commit:%s\n", formatValue(buildCommit))
}

func formatValue(buildData string) string {
	if buildData == "" {
		return "N/A"
	}
	return buildData
}

func main() {

	printBuildInfo()

	cfg := config.ParseOptions()
	a, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}

}
