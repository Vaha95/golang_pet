package handler

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

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

		id, err := setToStorage(storage, parsedURL.String())
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
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

func setToStorage(storage *repository.Storage, parsedURL string) (string, error) {
	for i := 0; i < 10; i++ {
		id := generateHash()
		if err := storage.Set(id, parsedURL); err != nil {
			if errors.Is(err, repository.ErrorShortURLKeyAlreadyExists) {
				continue
			}
			return "", fmt.Errorf("failed to save the short URL: %w", err)
		}

		return id, nil
	}

	return "", errors.New("ID generate is impossible")
}

func generateHash() string {
	rand.New((rand.NewSource(time.Now().UnixNano())))
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	return b.String()
}