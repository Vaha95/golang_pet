package handler

import "net/http"

func SaveUrl(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusCreated)
	res.Header().Set("content-type", "text/plain")
 }