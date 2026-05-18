package service

import (
	"context"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rafaeldepontes/voting-go/internal/middleware"
	"github.com/rafaeldepontes/voting-go/internal/poll"
	"github.com/rafaeldepontes/voting-go/internal/poll/model"
	"github.com/rafaeldepontes/voting-go/internal/utils"
	"github.com/rafaeldepontes/voting-go/internal/voting"
)

type service struct {
	mu          sync.RWMutex
	pr          poll.Repository
	subscribers map[string][]*websocket.Conn
}

func NewService(pr poll.Repository) voting.Service {
	return &service{
		pr:          pr,
		subscribers: make(map[string][]*websocket.Conn),
	}
}

// ListPolls implements [voting.Service].
func (s *service) ListPolls(ctx context.Context) []model.PollDto {
	polls, err := s.pr.ListPolls(ctx)
	if err != nil {
		log.Printf("[ERROR] couldn't list all polls: %v\n", err)
		return []model.PollDto{}
	}

	p := make([]model.PollDto, 0, len(polls))
	for i := range polls {
		p = append(p, model.PollDto{
			ID:       polls[i].ID,
			Text:     polls[i].Text,
			Duration: polls[i].Duration,
		})
	}
	return p
}

// CreatePoll implements [voting.Service].
func (s *service) CreatePoll(ctx context.Context, p model.PollReq) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := validatePollReq(p); err != nil {
		return "", err
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Printf("[ERROR] didn't create the uuid: %v\n", err)
		return "", utils.ErrGenericError
	}

	id := uuid.String()

	options := make([]model.Option, len(p.Options))
	for i := range p.Options {
		options[i] = model.Option{
			ID:    i + 1,
			Text:  p.Options[i],
			Votes: 0,
		}
	}

	ownerID := getUserID(ctx)
	poll := model.Poll{
		ID:        id,
		Text:      p.Name,
		Options:   options,
		Duration:  p.Duration,
		OwnerID:   ownerID,
		Votes:     make(map[string]struct{}),
		CreatedAt: time.Now(),
	}
	if err := s.pr.Insert(ctx, poll); err != nil {
		log.Printf("[ERROR] error inserting poll: %v\n", err)
		return "", utils.ErrGenericError
	}

	return id, nil
}

// RegisterVote implements [voting.Service].
func (s *service) RegisterVote(ctx context.Context, pollID string, optionID int) error {
	s.mu.Lock()
	poll, err := s.pr.FindPollByID(ctx, pollID)
	if err != nil {
		s.mu.Unlock()
		log.Printf("[ERROR] could not find poll by id %s because: %v\n", pollID, err)
		return utils.ErrPollNotFound
	}

	if poll.Duration > 0 && time.Since(poll.CreatedAt) > poll.Duration {
		s.mu.Unlock()
		return utils.ErrPollExpired
	}

	if _, has := poll.Votes[getUserID(ctx)]; has {
		s.mu.Unlock()
		return utils.ErrDuplicatedVote
	}

	found := false
	for i := range poll.Options {
		if poll.Options[i].ID == optionID {
			poll.Votes[getUserID(ctx)] = struct{}{}
			poll.Options[i].Votes++
			found = true
			break
		}
	}
	if err := s.pr.Update(ctx, poll); err != nil {
		log.Printf("[ERROR] error updating poll: %v\n", err)
		s.mu.Unlock()
		return utils.ErrGenericError
	}

	if !found {
		s.mu.Unlock()
		return utils.ErrOptionsNotFound
	}
	s.mu.Unlock()

	s.broadcast(ctx, pollID)
	return nil
}

// Subscribe implements [voting.Service].
func (s *service) Subscribe(ctx context.Context, pollID string, conn *websocket.Conn) error {
	s.mu.Lock()

	poll, err := s.pr.FindPollByID(ctx, pollID)
	if err != nil {
		s.mu.Unlock()
		log.Printf("[ERROR] could not find poll by id %s because: %v\n", pollID, err)
		return utils.ErrPollNotFound
	}

	s.subscribers[pollID] = append(s.subscribers[pollID], conn)
	s.mu.Unlock()

	return conn.WriteJSON(poll)
}

func (s *service) broadcast(ctx context.Context, pollID string) {
	s.mu.RLock()
	poll, err := s.pr.FindPollByID(ctx, pollID)
	if err != nil {
		s.mu.RUnlock()
		log.Printf("[ERROR] could not find poll by id %s because: %v\n", pollID, err)
		return
	}

	conns := s.subscribers[pollID]
	s.mu.RUnlock()

	var activeConns []*websocket.Conn
	for i := range conns {
		if err := conns[i].WriteJSON(poll); err != nil {
			continue
		}
		activeConns = append(activeConns, conns[i])
	}

	if len(activeConns) != len(conns) {
		s.mu.Lock()
		s.subscribers[pollID] = activeConns
		s.mu.Unlock()
	}
}

func getUserID(ctx context.Context) string {
	val := ctx.Value(middleware.UserID("userID"))
	if val == nil {
		return "anonymous"
	}

	switch v := val.(type) {
	case string:
		return v
	case int64:
		return strconv.FormatInt(v, 10)
	default:
		return "anonymous"
	}
}

func validatePollReq(p model.PollReq) error {
	if strings.TrimSpace(p.Name) == "" {
		return utils.ErrNameIsRequired
	}

	if len(p.Options) <= 1 {
		return utils.ErrOptionsIncorrectSize
	}

	if err := validateOptions(p.Options); err != nil {
		return err
	}

	return nil
}

func validateOptions(ops []string) error {
	for i := range ops {
		if strings.TrimSpace(ops[i]) == "" {
			return utils.ErrOptionsIsRequired
		}
	}
	return nil
}
