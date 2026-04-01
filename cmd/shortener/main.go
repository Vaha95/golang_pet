package main

import (
	"fmt"
	"log"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/config/db"
	"github.com/Vaha95/golang_pet/internal/handler"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	mv "github.com/Vaha95/golang_pet/internal/infrastructure/middleware"
)

func main() {
	mainConfig := config.GetMainConfig()

	dbService, err := db.InitDb(mainConfig) 
	defer dbService.Close()

	if err != nil {
		log.Fatal(
			fmt.Errorf("can`t init db: %w", err).Error(),
		)
	}

	isDBAllowed := dbService.Ping() == nil
	cfg := config.GetConfig(dbService, isDBAllowed, mainConfig)

	l, err := getLogger()
	if err != nil {
		log.Fatal(
			fmt.Errorf("can`t init logger: %w", err).Error(),
		)
	}

	e := echo.New()

	e.GET(`/:id`, handler.GetURLHandler(cfg, l))
	e.POST(`/`, handler.GetSaveURLHandler(cfg, l))
	e.POST(`/api/shorten`, handler.GetSaveURLShortenHandler(cfg, l))
	e.GET(`/ping`, handler.GetPingDBHandler(dbService))
	e.POST(`/api/shorten/batch`, handler.GetSaveURLBatchHandler(cfg, l))

	err = mv.AddMiddlewares(e, l)
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

func getLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("%w", mv.ErrorZapLoggerInitialize)
	}
	defer logger.Sync()

	return logger.Sugar(), nil
}