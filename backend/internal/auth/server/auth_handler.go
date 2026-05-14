package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/rafaeldepontes/voting-go/internal/auth"
	"github.com/rafaeldepontes/voting-go/internal/auth/model"
	jwt "github.com/rafaeldepontes/voting-go/internal/token"
	"github.com/rafaeldepontes/voting-go/internal/utils"
)

type handler struct {
	s          auth.Service
	jwtBuilder *jwt.JwtBuilder
}

func NewHandler(as auth.Service, jb *jwt.JwtBuilder) auth.Handler {
	return &handler{
		s:          as,
		jwtBuilder: jb,
	}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var userReq model.UserReq
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		log.Printf("[ERROR] failed to decode login request: %v\n", err)
		http.Error(w, utils.ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.s.Login(userReq)
	if err != nil {
		log.Printf("[ERROR] something went wrong log in: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenStr, _, err := h.jwtBuilder.GenerateToken(id, userReq.Email, 24*time.Hour)
	if err != nil {
		log.Printf("[ERROR] failed to generate token: %v\n", err)
		http.Error(w, utils.ErrFailedToGenerateToken.Error(), http.StatusInternalServerError)
		return
	}

	token := model.TokenResponse{
		ID:    id,
		Token: tokenStr,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var userReq model.UserReq
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		log.Printf("[ERROR] failed to decode register request: %v\n", err)
		http.Error(w, utils.ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	if err := h.s.Register(userReq); err != nil {
		log.Printf("[ERROR] couldn't register user: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
