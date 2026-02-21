package handler

import (
	"errors"
	"net/http"

	"github.com/Vaha95/golang_pet/internal/repository"
	"github.com/Vaha95/golang_pet/internal/service/save_url"
	"github.com/labstack/echo/v4"
)

func GetSaveURLShortenHandler(storage *repository.Storage, urlHost *string) (func(c echo.Context) error) {
	return func(c echo.Context) error {
		type Data struct {
			URI string `json:"url"`
		}
		var data Data

		if err := c.Bind(&data); err != nil {
			return c.String(http.StatusBadRequest, err.Error()) 
		}

		if urlHost == nil {
			link := `http://localhost:8080`
			urlHost = &link
		}

		path, err := saveurl.SaveUrl(data.URI, storage, *urlHost)
		if err != nil {
			switch errors.Is(err, saveurl.ErrorSaveToStorage) {
				case true:
					c.String(http.StatusInternalServerError, err.Error())							
				default:
					c.String(http.StatusBadRequest, err.Error())
			}	
		}

		return c.JSON(http.StatusCreated, path)
	}
}