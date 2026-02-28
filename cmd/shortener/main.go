package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Vaha95/golang_pet/internal/handler"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"github.com/labstack/echo/v4"

	mv "github.com/Vaha95/golang_pet/internal/infrastructure/middleware"
)

type Urls struct {
	listenHost string
	urlHost string
}

func main() {
	e := echo.New()
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(
			fmt.Errorf("can`t start logger: %w", err).Error(),
		)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	urls := getUrls()

	e.GET(`/:id`, handler.GetURLHandler())
	e.POST(`/`, handler.GetSaveURLHandler(&urls.urlHost))
	e.POST(`/api/shorten`, handler.GetSaveURLShortenHandler(&urls.urlHost))

	mv.AddMiddlewares(e, sugar)

	err = e.Start(urls.listenHost)
	if err != nil {
		log.Fatal(
			fmt.Errorf("can`t start Web server: %w", err).Error(),
		)
	}
}

// TODO вынести в config.go
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

	type baseURL struct {
		data string `env:"BASE_URL,required"`
	}
	var base baseURL
	err = env.Parse(&base)
	if err != nil || base.data == "" {
		urls.urlHost = *urlHostFlag
	} else {
		urls.urlHost = base.data
	}

	return urls
}