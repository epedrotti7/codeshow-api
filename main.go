package main

import (
	server "github.com/epedrotti7/codeshow-api/internal"
	router "github.com/epedrotti7/codeshow-api/internal/routers"
)

func main() {
	e := server.NewServer()
	router.StartRouter(e)
	server.StartServer(e)
}
