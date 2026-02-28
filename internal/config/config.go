package config

import (
	"flag"
	"os"
)

type Config struct {
	ListenHost string
	UrlHost string
}

func GetConfig() Config {
	listenHostENV := os.Getenv("SERVER_ADDRESS")
	urlHostENV := os.Getenv("BASE_URL")

	listenHostFlag := flag.String("a", `localhost:8080`, "Host for app")
	urlHostFlag := flag.String("b", `http://localhost:8080`, "Host for url")
	flag.Parse()

	listenHost := listenHostENV
	if listenHost == "" {
		listenHost = *listenHostFlag
	}

	urlHost := urlHostENV
	if urlHost == "" {
		urlHost = *urlHostFlag
	}

	return Config{
		ListenHost: listenHost,
		UrlHost: urlHost,
	}
}