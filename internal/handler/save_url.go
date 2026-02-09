package handler

import (
	"io"
	"net/http"
	"net/url"

	"github.com/Vaha95/golang_pet/internal/repository"
	echo "github.com/labstack/echo/v4"
)

func GetSaveURLHandler(storage *repository.Storage, urlHost *string) (func(c echo.Context) error) {
	return func(c echo.Context) error {
		req := c.Request()
		
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		inputURL := string(reqBody)
		parsedURL, err := url.ParseRequestURI(inputURL)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		id, err := storage.Set(parsedURL.String())
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		} else if (id == "") {
			return c.String(http.StatusInternalServerError, "Empty ID returned")
		}

		if urlHost == nil {
			link := `http://localhost:8080`
			urlHost = &link
		}

		path, err := url.JoinPath(*urlHost, id)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())			
		}

		return c.String(http.StatusCreated, path)
	}
}