package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	pm "github.com/rafaeldepontes/voting-go/internal/poll/model"
	"github.com/rafaeldepontes/voting-go/internal/utils"
	"github.com/rafaeldepontes/voting-go/internal/voting"
	vm "github.com/rafaeldepontes/voting-go/internal/voting/model"
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

// ListPolls implements [voting.Handler].
func (h *handler) ListPolls(w http.ResponseWriter, r *http.Request) {
	var p []pm.PollDto = h.s.ListPolls(r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

// CreatePoll implements [voting.Handler].
func (h *handler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	var req pm.PollReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.s.CreatePoll(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// CancelPoll implements [voting.Handler].
func (h *handler) CancelPoll(w http.ResponseWriter, r *http.Request) {
	pollID := r.PathValue("id")
	if pollID == "" {
		http.Error(w, utils.ErrPollIDMissing.Error(), http.StatusBadRequest)
		return
	}

	if err := h.s.CancelPoll(r.Context(), pollID); err != nil {
		if errors.Is(err, utils.ErrForbidden) {
			http.Error(w, err.Error(), http.StatusForbidden)
		}

		if errors.Is(err, utils.ErrPollNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		if errors.Is(err, utils.ErrGenericError) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterVote implements [voting.Handler].
func (h *handler) RegisterVote(w http.ResponseWriter, r *http.Request) {
	pollID := r.PathValue("id")
	if pollID == "" {
		http.Error(w, utils.ErrPollIDMissing.Error(), http.StatusBadRequest)
		return
	}

	var req vm.VoteReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.s.RegisterVote(r.Context(), pollID, req.OptionID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleWS implements [voting.Handler].
func (h *handler) HandleWS(w http.ResponseWriter, r *http.Request) {
	pollID := r.PathValue("id")
	if pollID == "" {
		http.Error(w, utils.ErrPollIDMissing.Error(), http.StatusBadRequest)
		return
	}

	var conn *websocket.Conn
	conn, err := h.u.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR] websocket upgrade failed: %v\n", err)
		http.Error(w, utils.ErrGenericError.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.s.Subscribe(r.Context(), pollID, conn); err != nil {
		conn.Close()
		log.Printf("[ERROR] subscription failed: %v\n", err)
		http.Error(w, utils.ErrGenericError.Error(), http.StatusInternalServerError)
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
