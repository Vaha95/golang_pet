package main

import (
	"net/http"

	"github.com/Vaha95/golang_pet/internal/handler"
	"github.com/gorilla/mux"
)

func getEndpoints() {
	mux := mux.NewRouter()
	mux.HandleFunc(`/`, handler.SaveUrl)
	mux.HandleFunc(`/{id}`, handler.GetUrl)
	listen(`:8080`, mux)
}

func listen(addr string, handler http.Handler) {
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		panic(err)
	}
}