package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Vaha95/golang_pet/internal/handler"
	"github.com/Vaha95/golang_pet/internal/repository"
	"github.com/caarlos0/env/v6"
	echo "github.com/labstack/echo/v4"
)

type Urls struct {
	listenHost string
	urlHost string
}

func main() {
	e := echo.New()

	urls := getUrls()

	storage := repository.NewStorage()

	e.GET(`/:id`, handler.GetURLHandler(storage))
	e.POST(`/`, handler.GetSaveURLHandler(storage, &urls.urlHost))

	err := e.Start(*&urls.listenHost)
	if err != nil {
		log.Fatal(
			fmt.Errorf("can`t start Web server: %w", err).Error(),
		)
	}
}

func getUrls() Urls {
	var urls Urls

	listenHostFlag := flag.String("a", `localhost:8080`, "Host for app")
	urlHostFlag := flag.String("b", `http://localhost:8080`, "Host for url")
	flag.Parse()

	type serverAddr struct {
		data string `env:"SERVER_ADDRESS,required"`
	}
	var addr serverAddr
	err := env.Parse(&addr)
	if err != nil || addr.data == ""{
		urls.listenHost = *listenHostFlag
	} else {
		urls.listenHost = addr.data
	}

	type baseUrl struct {
		data string `env:"BASE_URL,required"`
	}
	var base baseUrl
	err = env.Parse(&base)
	if err != nil || base.data == "" {
		urls.urlHost = *urlHostFlag
	} else {
		urls.urlHost = base.data
	}

	return urls
}