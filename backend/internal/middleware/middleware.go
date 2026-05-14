package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/rafaeldepontes/voting-go/internal/auth"
	jwt "github.com/rafaeldepontes/voting-go/internal/token"
	"github.com/rafaeldepontes/voting-go/internal/utils"
	"github.com/rafaeldepontes/voting-go/internal/voting"
)

const (
	TokenPrefix = "Bearer "
)

type Middleware struct {
	JwtBuilder *jwt.JwtBuilder
}

func NewMiddleware(sk string) Middleware {
	return Middleware{
		JwtBuilder: jwt.NewJwtBuilder(sk),
	}
}

func (m *Middleware) AuthFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		dirtToken := r.Header.Get("Authorization")

		if dirtToken != "" {
			if !strings.HasPrefix(dirtToken, TokenPrefix) {
				log.Printf("[DEBUG] Authorization header missing prefix: %s\n", dirtToken)
				http.Error(w, utils.ErrInvalidToken.Error(), http.StatusForbidden)
				return
			}
			token = cleanToken(dirtToken)
		} else {
			// Check query parameter for WebSocket connections
			token = r.URL.Query().Get("token")
		}

		if token == "" {
			log.Println("[DEBUG] Token is missing (neither in header nor query param)")
			http.Error(w, utils.ErrInvalidToken.Error(), http.StatusForbidden)
			return
		}

		_, err := m.JwtBuilder.VerifyToken(token)
		if err != nil {
			log.Printf("[DEBUG] token verification failed: %v, token length: %d\n", err, len(token))
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewHandler(frontend string, m Middleware, han voting.Handler, auth auth.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", auth.Register)
	mux.HandleFunc("POST /login", auth.Login)

	// Polls
	mux.Handle(
		"POST /polls", m.AuthFilter(
			http.HandlerFunc(han.CreatePoll),
		),
	)
	mux.Handle(
		"GET /ws/polls/{id}", m.AuthFilter(
			http.HandlerFunc(han.HandleWS),
		),
	)
	mux.Handle(
		"GET /polls", m.AuthFilter(
			http.HandlerFunc(han.ListPolls),
		),
	)

	// Votes
	mux.Handle(
		"POST /polls/{id}/vote", m.AuthFilter(
			http.HandlerFunc(han.RegisterVote),
		),
	)

	return mux
}
