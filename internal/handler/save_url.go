package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/service/save_url"
	"github.com/labstack/echo/v4"
)

func GetSaveURLHandler(cfg config.StorageConfig) (func(c echo.Context) error) {
	return func(c echo.Context) error {
		req := c.Request()
		
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		inputURL := string(reqBody)

		path, err := saveurl.SaveURL(cfg, inputURL)
		if err != nil {
			 if (errors.Is(err, saveurl.ErrorUrlAlreadyExists)) {
				return c.JSON(http.StatusConflict, path)
			} else if errors.Is(err, saveurl.ErrorSaveToStorage) {
				return c.JSON(http.StatusInternalServerError, err.Error())
			} else {
				return c.JSON(http.StatusBadRequest, err.Error())
			}
		}

		return c.String(http.StatusCreated, path)
	}
}