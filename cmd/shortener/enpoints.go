package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	// "github.com/Vaha95/golang_pet/internal/handler"
	"github.com/gorilla/mux"
)

var data map[string]string
type Storage struct {
	data sync.Map
}

func (s *Storage) Set(k string, v interface{}) {
	s.data.Store(k, v)
}

func (s *Storage) Get(k string) (interface{}, bool) {
	return s.data.Load(k)
}

func NewStorage () *Storage {
	return &Storage{}
}

var storage = NewStorage()

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
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests allowed!", http.StatusMethodNotAllowed)

		return
	}

	err := req.ParseForm()
	if err != nil {
		res.Write([]byte(err.Error()))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	id := ""
	for _, v := range req.Form {
		id = generateId()
		storage.Set(id, v[0])
	}

	res.WriteHeader(http.StatusCreated)
	res.Header().Set("content-type", "text/plain")

	res.Write(fmt.Appendf(nil, "http://localhost:8080/%s", id))
 }

func getUrl(res http.ResponseWriter, req *http.Request) {
	// if req.Method != http.MethodGet {
	// 	http.Error(res, "Only GET requests allowed!", http.StatusMethodNotAllowed)

	// 	return
	// }

	vars := mux.Vars(req)
    id := vars["id"]

	val, ok := storage.Get(id)
	if !ok {		
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(""))
	}

	res.WriteHeader(http.StatusTemporaryRedirect)
	res.Header().Set("content-type", "text/plain")
	res.Header().Set("location", val.(string))

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