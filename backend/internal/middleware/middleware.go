package middleware

import (
	"context"
	"log"
	"net"
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

type UserID string

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
			token = r.URL.Query().Get("token")
		}

		// If the token is not present in the request... The user should be treated as anonymous
		// if token == "" {
		// 	log.Println("[DEBUG] Token is missing (neither in header nor query param)")
		// 	http.Error(w, utils.ErrInvalidToken.Error(), http.StatusForbidden)
		// 	return
		// }

		ip := GetIP(r)
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserID("userID"), ip)

		if token != "" {
			user, err := m.JwtBuilder.VerifyToken(token)
			if err != nil {
				log.Printf("[DEBUG] token verification failed: %v, token length: %d\n", err, len(token))
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx = context.WithValue(ctx, UserID("userID"), user.ID)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
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

func GetIP(r *http.Request) string {
	for _, header := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		if ip := r.Header.Get(header); ip != "" {
			return strings.Split(ip, ",")[0]
		}
	}

	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}
