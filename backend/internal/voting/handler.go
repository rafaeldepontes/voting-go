package voting

import "net/http"

type Handler interface {
	ListPolls(w http.ResponseWriter, r *http.Request)
	CreatePoll(w http.ResponseWriter, r *http.Request)
	CancelPoll(w http.ResponseWriter, r *http.Request)
	RegisterVote(w http.ResponseWriter, r *http.Request)
	HandleWS(w http.ResponseWriter, r *http.Request)
}
