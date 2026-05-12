package servers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rafaeldepontes/voting-go/internal/poll/model"
	"github.com/rafaeldepontes/voting-go/internal/utils"
	"github.com/rafaeldepontes/voting-go/internal/voting"
)

type handler struct {
	u websocket.Upgrader
	s voting.Service
}

func NewHandler(u websocket.Upgrader, s voting.Service) voting.Handler {
	return &handler{
		u: u,
		s: s,
	}
}

// CreatePoll implements [voting.Handler].
func (h *handler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	var req model.PollReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.s.CreatePoll(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// RegisterVote implements [voting.Handler].
func (h *handler) RegisterVote(w http.ResponseWriter, r *http.Request) {
	pollID := r.PathValue("id")
	if pollID == "" {
		http.Error(w, utils.PollIDMissing.Error(), http.StatusBadRequest)
		return
	}

	var req model.VoteReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.s.RegisterVote(pollID, req.OptionID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleWS implements [voting.Handler].
func (h *handler) HandleWS(w http.ResponseWriter, r *http.Request) {
	pollID := r.PathValue("id")
	if pollID == "" {
		http.Error(w, utils.PollIDMissing.Error(), http.StatusBadRequest)
		return
	}

	var conn *websocket.Conn
	conn, err := h.u.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR] websocket upgrade failed: %v\n", err)
		http.Error(w, utils.GenericError.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.s.Subscribe(pollID, conn); err != nil {
		conn.Close()
		log.Printf("[ERROR] subscription failed: %v\n", err)
		http.Error(w, utils.GenericError.Error(), http.StatusInternalServerError)
		return
	}

	go func() {
		defer conn.Close()
		for {
			if _, _, err := conn.NextReader(); err != nil {
				break
			}
		}
	}()
}
