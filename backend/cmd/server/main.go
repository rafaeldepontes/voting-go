package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	ar "github.com/rafaeldepontes/voting-go/internal/auth/repository"
	asr "github.com/rafaeldepontes/voting-go/internal/auth/server"
	asv "github.com/rafaeldepontes/voting-go/internal/auth/service"
	"github.com/rafaeldepontes/voting-go/internal/middleware"
	pr "github.com/rafaeldepontes/voting-go/internal/poll/repository"
	vsr "github.com/rafaeldepontes/voting-go/internal/voting/server"
	vsv "github.com/rafaeldepontes/voting-go/internal/voting/service"
	"github.com/rafaeldepontes/voting-go/pkg/database/postgres"
)

const (
	ReadBufferSize  = 4096
	WriteBufferSize = 4096
)

var u = websocket.Upgrader{
	WriteBufferSize: WriteBufferSize,
	ReadBufferSize:  ReadBufferSize,
}

func init() {
	_ = godotenv.Load(".env")
}

func main() {
	port := os.Getenv("PORT")
	origin := os.Getenv("ORIGIN_URL")
	frontend := os.Getenv("FRONTEND_URL")
	if port == "" {
		port = "8080"
	}

	u.CheckOrigin = func(r *http.Request) bool {
		return r.Host == origin
	}
	// u.CheckOrigin = func(r *http.Request) bool { return true }

	postgres.GetDb()
	defer postgres.Close()

	if err := postgres.RunMigrations(); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	r := pr.NewRepository()
	ss := vsv.NewService(r)
	han := vsr.NewHandler(u, ss)

	ar := ar.NewRepository()
	as := asv.NewService(ar)
	m := middleware.NewMiddleware(os.Getenv("SECRET_KEY"))
	auth := asr.NewHandler(as, m.JwtBuilder)

	mux := middleware.NewHandler(frontend, m, han, auth)
	corsMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", frontend)
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
