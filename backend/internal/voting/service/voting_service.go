package service

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rafaeldepontes/voting-go/internal/poll/model"
	"github.com/rafaeldepontes/voting-go/internal/utils"
	"github.com/rafaeldepontes/voting-go/internal/voting"
)

type service struct {
	mu          sync.RWMutex
	polls       map[string]*model.Poll
	subscribers map[string][]*websocket.Conn
	pollCounter int
}

func NewService() voting.Service {
	return &service{
		polls:       make(map[string]*model.Poll),
		subscribers: make(map[string][]*websocket.Conn),
	}
}

func (s *service) ListPolls() []model.PollDto {
	p := make([]model.PollDto, 0, len(s.polls))
	for i := range s.polls {
		p = append(p, model.PollDto{
			ID:   s.polls[i].ID,
			Text: s.polls[i].Text,
		})
	}
	return p
}

func (s *service) CreatePoll(p model.PollReq) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pollCounter++
	id := fmt.Sprintf("%d", s.pollCounter)

	options := make([]model.Option, len(p.Options))
	for i := range p.Options {
		options[i] = model.Option{
			ID:    i + 1,
			Text:  p.Options[i],
			Votes: 0,
		}
	}

	s.polls[id] = &model.Poll{
		ID:      id,
		Text:    p.Name,
		Options: options,
	}

	return id, nil
}

// RegisterVote implements [voting.Service].
func (s *service) RegisterVote(pollID string, optionID int) error {
	s.mu.Lock()
	poll, ok := s.polls[pollID]
	if !ok {
		s.mu.Unlock()
		return utils.PollNotFound
	}

	found := false
	for i := range poll.Options {
		if poll.Options[i].ID == optionID {
			poll.Options[i].Votes++
			found = true
			break
		}
	}
	s.mu.Unlock()

	if !found {
		return utils.OptionsNotFound
	}

	s.broadcast(pollID)
	return nil
}

func (s *service) Subscribe(pollID string, conn *websocket.Conn) error {
	s.mu.Lock()

	poll, ok := s.polls[pollID]
	if !ok {
		s.mu.Unlock()
		return utils.PollNotFound
	}

	s.subscribers[pollID] = append(s.subscribers[pollID], conn)
	s.mu.Unlock()

	return conn.WriteJSON(poll)
}

func (s *service) broadcast(pollID string) {
	s.mu.RLock()
	poll := s.polls[pollID]
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
