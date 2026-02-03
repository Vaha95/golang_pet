package main

import (
	"net/http"

	"github.com/Vaha95/golang_pet/internal/handler"
	"github.com/gorilla/mux"
)

func getEndpoints() {
	mux := mux.NewRouter()

	mux.HandleFunc(`/{id}`, handler.GetGetURLHandler(&Urls)).Methods(http.MethodGet)
	mux.HandleFunc(`/`, handler.GetSaveURLHandler(&Urls)).Methods(http.MethodPost)

	listen(`localhost:8080`, mux)
}

func listen(addr string, handler http.Handler) {
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		panic(err)
	}
}