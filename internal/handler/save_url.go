package handler

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetSaveURLHandler(data *map[string]string) (func(res http.ResponseWriter, req *http.Request)) {
	saveURL := func (res http.ResponseWriter, req *http.Request)  {
		if req.Method != http.MethodPost {
			http.Error(res, "Only POST requests allowed!", http.StatusMethodNotAllowed)

			return
		}

		defer req.Body.Close()
		
		reqBody, err := io.ReadAll(req.Body)

		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)

			return				
		}

		inputURL := string(reqBody)
		u, err := url.ParseRequestURI(inputURL)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)

			return
		}

		id := generateID()
		(*data)[id] = u.String()

		res.WriteHeader(http.StatusCreated)
		res.Header().Set("Content-type", "text/plain")

		res.Write(fmt.Appendf(nil, "http://localhost:8080/%s", id))
	}

	return saveURL
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