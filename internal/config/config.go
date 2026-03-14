package config

import (
	"flag"
	"os"
)

type Config struct {
	ListenHost string
	URLHost string
	FilePath string
	DbDSN string
}

func GetConfig() Config {
	listenHostENV := os.Getenv("SERVER_ADDRESS")
	urlHostENV := os.Getenv("BASE_URL")
	filePathENV := os.Getenv("FILE_STORAGE_PATH")
	dsnENV := os.Getenv("DATABASE_DSN")

	listenHostFlag := flag.String("a", `localhost:8080`, "Host for app")
	urlHostFlag := flag.String("b", `http://localhost:8080`, "Host for url")
	filePathFlag := flag.String("f", ``, "Storage file path")
	dsnFlag := flag.String("d", ``, "Database dsn")
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

	dsn := dsnENV
	if dsn == "" {
		dsn = *dsnFlag
	}

	return Config{
		ListenHost: listenHost,
		URLHost: urlHost,
		FilePath: filePath,
		DbDSN: dsn,
	}
}