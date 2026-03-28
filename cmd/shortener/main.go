package main

import (
	"fmt"
	"log"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/config/db"
	"github.com/Vaha95/golang_pet/internal/handler"
	"github.com/labstack/echo/v4"

	mv "github.com/Vaha95/golang_pet/internal/infrastructure/middleware"
)

func main() {
	mainConfig := config.GetMainConfig()

	dbService := db.InitDb(mainConfig) 
	defer dbService.Close()

	isDBAllowed := dbService.Ping() == nil
	cfg := config.GetConfig(dbService, isDBAllowed, mainConfig)

	e := echo.New()

	e.GET(`/:id`, handler.GetURLHandler(cfg.Config))
	e.POST(`/`, handler.GetSaveURLHandler(cfg))
	e.POST(`/api/shorten`, handler.GetSaveURLShortenHandler(cfg))
	e.GET(`/ping`, handler.GetPingDBHandler(dbService))
	e.POST(`/api/shorten/batch`, handler.GetSaveURLBatchHandler(cfg))

	err := mv.AddMiddlewares(e)
	if err != nil {
		log.Fatal(
			fmt.Errorf("can`t start Web server: %w", err).Error(),
		)
	}

	err = e.Start(cfg.Config.ListenHost)
	if err != nil {
		log.Fatal(
			fmt.Errorf("can`t start Web server: %w", err).Error(),
		)
	}
}