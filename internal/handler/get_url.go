package handler

import (
	"net/http"

	// "github.com/gorilla/mux"
)

func GetUrl(res http.ResponseWriter, req *http.Request) {    
	// vars := mux.Vars(req)
    // id := vars["id"]
	res.WriteHeader(http.StatusTemporaryRedirect)
	res.Header().Set("content-type", "text/plain")
	res.Header().Set("location", "text/plain")
}