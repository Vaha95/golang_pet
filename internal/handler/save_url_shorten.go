package handler

import (
	"errors"
	"net/http"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/service/save_url"
	"github.com/labstack/echo/v4"
)

func GetSaveURLShortenHandler(cfg config.StorageConfig) (func(c echo.Context) error) {
	return func(c echo.Context) error {
		type APIReqiest struct {
			URI string `json:"url"`
		}
		var data APIReqiest

		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error()) 
		}

		path, err := saveurl.SaveURL(cfg, data.URI)
		if err != nil {
			if (errors.Is(err, saveurl.ErrorUrlAlreadyExists)) {
				return c.String(http.StatusConflict, path)
			} else if errors.Is(err, saveurl.ErrorSaveToStorage) {
				return c.JSON(http.StatusInternalServerError, err.Error())
			} else {
				return c.JSON(http.StatusBadRequest, err.Error())
			}
		}

		type APIResponse struct {
			Result string `json:"result"`
		}

		return c.JSON(http.StatusCreated, APIResponse{Result: path})
	}
}