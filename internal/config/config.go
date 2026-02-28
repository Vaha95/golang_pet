package config

import (
	"flag"
	"os"
)

type Config struct {
	ListenHost string
	URLHost string
	FilePath string
}

func GetConfig() Config {
	listenHostENV := os.Getenv("SERVER_ADDRESS")
	urlHostENV := os.Getenv("BASE_URL")
	filePathENV := os.Getenv("FILE_STORAGE_PATH")

	listenHostFlag := flag.String("a", `localhost:8080`, "Host for app")
	urlHostFlag := flag.String("b", `http://localhost:8080`, "Host for url")
	filePathFlag := flag.String("f", ``, "Storage file path")
	flag.Parse()

	listenHost := listenHostENV
	if listenHost == "" {
		listenHost = *listenHostFlag
	}

	urlHost := urlHostENV
	if urlHost == "" {
		urlHost = *urlHostFlag
	}

	filePath := filePathENV
	if filePath == "" {
		filePath = *filePathFlag
	}

	return Config{
		ListenHost: listenHost,
		URLHost: urlHost,
		FilePath: filePath,
	}
}