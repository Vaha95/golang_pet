package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	// "github.com/Vaha95/golang_pet/internal/handler"
	"github.com/gorilla/mux"
)

var data map[string]string

func getEndpoints() {
	mux := mux.NewRouter()

	mux.HandleFunc(`/`, saveUrl)
	mux.HandleFunc(`/{id}`, getUrl)

	listen(`:8080`, mux)
}

func listen(addr string, handler http.Handler) {
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		panic(err)
	}
}

func saveUrl(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		res.Write([]byte(err.Error()))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	id := ""
	for _, v := range req.Form {
		id = generateId()
		data[id] = v[0]
	}

	res.WriteHeader(http.StatusCreated)
	res.Header().Set("content-type", "text/plain")

	res.Write([]byte(fmt.Sprintf("localhost:8080/%s", id)))
 }

func getUrl(res http.ResponseWriter, req *http.Request) {    
	vars := mux.Vars(req)
    id := vars["id"]

	res.WriteHeader(http.StatusTemporaryRedirect)
	res.Header().Set("content-type", "text/plain")
	res.Header().Set("location", data[id])

	res.Write([]byte(""))
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