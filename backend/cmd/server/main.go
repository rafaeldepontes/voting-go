package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/voting-go/internal/poll/repository"
	server "github.com/rafaeldepontes/voting-go/internal/voting/server"
	"github.com/rafaeldepontes/voting-go/internal/voting/service"
	"github.com/rafaeldepontes/voting-go/pkg/cache/redis"
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

func init() {
	_ = godotenv.Load(".env")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := repository.NewRepository(redis.GetCache())
	ss := service.NewService(r)
	h := server.NewHandler(u, ss)
	mux := server.MapRoutesPoll(h)

	corsMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		mux.ServeHTTP(w, r)
	})

	log.Printf("Application running on port %s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, corsMux))
}
