package servers

import (
	"net/http"

	"github.com/rafaeldepontes/voting-go/internal/voting"
)

func MapRoutesPoll(h voting.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /polls", h.CreatePoll)
	mux.HandleFunc("POST /polls/{id}/vote", h.RegisterVote)
	mux.HandleFunc("GET /ws/polls/{id}", h.HandleWS)
	mux.HandleFunc("GET /polls", h.ListPolls)

	return mux
}
