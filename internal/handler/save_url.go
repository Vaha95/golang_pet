package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	saveurl "github.com/Vaha95/golang_pet/internal/service/save_url"
)

func GetSaveURLHandler(cfg config.StorageConfig, l *zap.SugaredLogger) func(c echo.Context) error {
	return func(c echo.Context) error {
		req := c.Request()

		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		inputURL := string(reqBody)

		path, err := saveurl.SaveURL(cfg, inputURL)
		if err != nil {
			if errors.Is(err, saveurl.ErrorUrlAlreadyExists) {
				return c.String(http.StatusConflict, path)
			} else if errors.Is(err, saveurl.ErrorSaveToStorage) {
				l.Errorf("Default url save error: %w", err)

				return c.NoContent(http.StatusInternalServerError)
			} else {
				return c.JSON(http.StatusBadRequest, err.Error())
			}
		}

		return c.String(http.StatusCreated, path)
	}
}
