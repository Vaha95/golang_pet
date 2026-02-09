package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Vaha95/golang_pet/internal/handler"
	"github.com/Vaha95/golang_pet/internal/repository"
	echo "github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	listenHost := flag.String("a", `localhost:8080`, "Host for app")
	urlHost := flag.String("b", `http://localhost:8080`, "Host for url")
	flag.Parse()

	storage := repository.NewStorage()

	e.GET(`/:id`, handler.GetURLHandler(storage))
	e.POST(`/`, handler.GetSaveURLHandler(storage, urlHost))

	error := listen(e, *listenHost)
	if error != nil {
		log.Fatal(
			fmt.Errorf("Can`t start Web server: %w", error).Error(),
		)
	}
}

func listen(e *echo.Echo, addr string) error {
	err := e.Start(addr)

	return err
}
