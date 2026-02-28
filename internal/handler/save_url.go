package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/service/save_url"
	"github.com/labstack/echo/v4"
)

func GetSaveURLHandler(cfg config.Config) (func(c echo.Context) error) {
	return func(c echo.Context) error {
		req := c.Request()
		
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		inputURL := string(reqBody)

		path, err := saveurl.SaveURL(cfg, inputURL)
		if err != nil {
			switch errors.Is(err, saveurl.ErrorSaveToStorage) {
				case true:
					c.String(http.StatusInternalServerError, err.Error())							
				default:
					c.String(http.StatusBadRequest, err.Error())
			}	
		}

		return c.String(http.StatusCreated, path)
	}
}