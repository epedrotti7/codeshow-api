package server

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewServer() *echo.Echo {
	e := echo.New()
	return e
}

func StartServer(e *echo.Echo) {
	if err := e.Start(":3000"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
