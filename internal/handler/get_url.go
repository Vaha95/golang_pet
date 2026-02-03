package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func GetGetURLHandler(data *map[string]string) (func(res http.ResponseWriter, req *http.Request)) {
	saveURL := func (res http.ResponseWriter, req *http.Request)  {
	vars := mux.Vars(req)
    id := vars["id"]

	val, ok := (*data)[id]
	if !ok || val == "" {		
		res.WriteHeader(http.StatusNotFound)

		return
	}

	http.Redirect(res, req, val, http.StatusTemporaryRedirect)
	}

	return saveURL
}