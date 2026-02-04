package main

import (
	"github.com/Vaha95/golang_pet/internal/handler"
	echo "github.com/labstack/echo/v4"
)

func getEndpoints() {
	e := echo.New()

	e.GET(`/:id`, handler.GetGetURLHandler(&Urls))
	e.POST(`/`, handler.GetSaveURLHandler(&Urls))

	listen(e, `localhost:8080`)
}

func listen(e *echo.Echo, addr string) {
	err := e.Start(addr)
	if err != nil {
		panic(err)
	}
}