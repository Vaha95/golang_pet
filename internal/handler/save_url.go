package handler

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	echo "github.com/labstack/echo/v4"
)

func GetSaveURLHandler(data *map[string]string, urlHost *string) (func(c echo.Context) error) {
	return func(c echo.Context) error {
		req := c.Request()
		defer req.Body.Close()
		
		reqBody, err := io.ReadAll(req.Body)

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		inputURL := string(reqBody)
		u, err := url.ParseRequestURI(inputURL)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		id := generateID()
		(*data)[id] = u.String()

		if urlHost == nil {
			link := `http://localhost:8080`
			urlHost = &link
		}

		return c.String(http.StatusCreated, fmt.Sprintf("%s/%s", *urlHost, id))
	}
}

func generateID() (string) {
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