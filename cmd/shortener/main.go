package main

import (
	"fmt"
	"log"

	"github.com/Vaha95/golang_pet/internal/handler"
	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/labstack/echo/v4"

	mv "github.com/Vaha95/golang_pet/internal/infrastructure/middleware"
)

func main() {
	e := echo.New()

	cfg := config.GetConfig()

	e.GET(`/:id`, handler.GetURLHandler(cfg))
	e.POST(`/`, handler.GetSaveURLHandler(cfg))
	e.POST(`/api/shorten`, handler.GetSaveURLShortenHandler(cfg))

	err := mv.AddMiddlewares(e)
	if err != nil {
		log.Fatal(
			fmt.Errorf("can`t start Web server: %w", err).Error(),
		)
	}

	err = e.Start(cfg.ListenHost)
	if err != nil {
		log.Fatal(
			fmt.Errorf("can`t start Web server: %w", err).Error(),
		)
	}
}