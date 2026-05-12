package voting

import "net/http"

type Handler interface {
	CreatePoll(w http.ResponseWriter, r *http.Request)
	RegisterVote(w http.ResponseWriter, r *http.Request)
	HandleWS(w http.ResponseWriter, r *http.Request)
}
