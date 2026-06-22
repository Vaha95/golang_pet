package handler

import (
	"fmt"
	"net/http"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

func GetDeleteURLHandler(cfg config.StorageConfig, ch chan DTO.DeleteBatch, l *zap.SugaredLogger) func(c *echo.Context) error {
	return func(c *echo.Context) error {
		var req []string

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if len(req) <= 0 {
			return c.JSON(http.StatusBadRequest, fmt.Errorf("Request data is empty!"))
		}

		msg := DTO.DeleteBatch{
			UserId: *cfg.GetUserId(),
			Shorts: req,
		}

		ch <- msg

		return c.JSON(http.StatusAccepted, nil)
	}
}
