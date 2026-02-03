package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var urls map[string]string

func init() {
	urls = make(map[string]string)
}

func main() {
	mux := mux.NewRouter()

	mux.HandleFunc(`/{id}`, getUrl).Methods(http.MethodGet)
	mux.HandleFunc(`/`, saveUrl).Methods(http.MethodPost)

	listen(`localhost:8080`, mux)
}

func listen(addr string, handler http.Handler) {
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		panic(err)
	}
}

func saveUrl(res http.ResponseWriter, req *http.Request) {
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

	inputUrl := string(reqBody)
	u, err := url.ParseRequestURI(inputUrl)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}

	id := generateId()
	urls[id] = u.String()

	res.WriteHeader(http.StatusCreated)
	res.Header().Set("Content-type", "text/plain")

	res.Write(fmt.Appendf(nil, "http://localhost:8080/%s", id))
 }

func getUrl(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
    id := vars["id"]

	val, ok := urls[id]
	if !ok || val == "" {		
		res.WriteHeader(http.StatusNotFound)
		data, _ := json.Marshal(vars)
		res.Write(data)

		return
	}

	http.Redirect(res, req, val, http.StatusTemporaryRedirect)
}

func generateId() (string) {
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