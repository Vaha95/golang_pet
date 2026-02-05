package main

import (
	"flag"

	"github.com/Vaha95/golang_pet/internal/handler"
	echo "github.com/labstack/echo/v4"
)

var Urls map[string]string

func init() {
	Urls = make(map[string]string)
}

func main() {
	e := echo.New()

	listenHost := flag.String("a", `localhost:8080`, "Host for app")
	urlHost := flag.String("b", `http://localhost:8080`, "Host for url")
	flag.Parse()

	e.GET(`/:id`, handler.GetGetURLHandler(&Urls))
	e.POST(`/`, handler.GetSaveURLHandler(&Urls, urlHost))

	listen(e, *listenHost)
}

func listen(e *echo.Echo, addr string) {
	err := e.Start(addr)
	if err != nil {
		panic(err)
	}
}
