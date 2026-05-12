package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	server "github.com/rafaeldepontes/voting-go/internal/voting/server"
	"github.com/rafaeldepontes/voting-go/internal/voting/service"
)

const (
	ReadBufferSize  = 4096
	WriteBufferSize = 4096
)

var u = websocket.Upgrader{
	WriteBufferSize: WriteBufferSize,
	ReadBufferSize:  ReadBufferSize,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ss := service.NewService()
	h := server.NewHandler(u, ss)
	mux := server.MapRoutesPoll(h)

	log.Printf("Application running on port %s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, mux))
}
