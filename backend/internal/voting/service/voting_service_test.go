package service

import (
	"context"
	"testing"
	"time"

	"github.com/rafaeldepontes/voting-go/internal/middleware"
	"github.com/rafaeldepontes/voting-go/internal/poll/model"
	"github.com/rafaeldepontes/voting-go/internal/utils"
)

type mockRepository struct {
	polls map[string]model.Poll
	err   error
}

func (m *mockRepository) Insert(ctx context.Context, poll model.Poll) error {
	if m.err != nil {
		return m.err
	}
	m.polls[poll.ID] = poll
	return nil
}

func (m *mockRepository) FindPollByID(ctx context.Context, id string) (model.Poll, error) {
	if m.err != nil {
		return model.Poll{}, m.err
	}
	p, ok := m.polls[id]
	if !ok {
		return model.Poll{}, utils.ErrPollNotFound
	}
	return p, nil
}

func (m *mockRepository) ListPolls(ctx context.Context) ([]model.Poll, error) {
	if m.err != nil {
		return nil, m.err
	}
	res := make([]model.Poll, 0, len(m.polls))
	for _, p := range m.polls {
		res = append(res, p)
	}
	return res, nil
}

func (m *mockRepository) Update(ctx context.Context, poll model.Poll) error {
	if m.err != nil {
		return m.err
	}
	m.polls[poll.ID] = poll
	return nil
}

func (m *mockRepository) Remove(ctx context.Context, id string) error {
	if m.err != nil {
		return m.err
	}
	if _, ok := m.polls[id]; !ok {
		return utils.ErrPollNotFound
	}
	delete(m.polls, id)
	return nil
}

func TestCreatePoll(t *testing.T) {
	repo := &mockRepository{polls: make(map[string]model.Poll)}
	s := NewService(repo)

	tests := []struct {
		name    string
		req     model.PollReq
		wantErr error
	}{
		{
			name: "Success",
			req: model.PollReq{
				Name:    "Favorite Color",
				Options: []string{"Red", "Blue", "Green"},
			},
			wantErr: nil,
		},
		{
			name: "Missing Name",
			req: model.PollReq{
				Name:    "",
				Options: []string{"Red", "Blue"},
			},
			wantErr: utils.ErrNameIsRequired,
		},
		{
			name: "Insufficient Options",
			req: model.PollReq{
				Name:    "Test",
				Options: []string{"Only One"},
			},
			wantErr: utils.ErrOptionsIncorrectSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := s.CreatePoll(context.Background(), tt.req)
			if err != tt.wantErr {
				t.Errorf("CreatePoll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil && id == "" {
				t.Error("CreatePoll() returned empty ID")
			}
		})
	}
}

func TestCancelPoll(t *testing.T) {
	ownerID := "user123"
	pollID := "poll123"

	tests := []struct {
		name    string
		pollID  string
		userID  string
		wantErr error
	}{
		{
			name:    "Success",
			pollID:  pollID,
			userID:  ownerID,
			wantErr: nil,
		},
		{
			name:    "Not Found",
			pollID:  "wrong-id",
			userID:  ownerID,
			wantErr: utils.ErrPollNotFound,
		},
		{
			name:    "Forbidden",
			pollID:  pollID,
			userID:  "other-user",
			wantErr: utils.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{polls: make(map[string]model.Poll)}
			repo.polls[pollID] = model.Poll{
				ID:      pollID,
				OwnerID: ownerID,
			}
			s := NewService(repo)

			ctx := context.WithValue(context.Background(), middleware.UserInfo("userID"), tt.userID)
			err := s.CancelPoll(ctx, tt.pollID)
			if err != tt.wantErr {
				t.Errorf("CancelPoll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterVote(t *testing.T) {
	pollID := "poll123"

	tests := []struct {
		name     string
		pollID   string
		userID   string
		optionID int
		setup    func(repo *mockRepository)
		wantErr  error
	}{
		{
			name:     "Success",
			pollID:   pollID,
			userID:   "user1",
			optionID: 1,
			wantErr:  nil,
		},
		{
			name:     "Duplicated Vote",
			pollID:   pollID,
			userID:   "user1",
			optionID: 2,
			setup: func(repo *mockRepository) {
				p := repo.polls[pollID]
				p.Votes["user1"] = struct{}{}
				repo.polls[pollID] = p
			},
			wantErr: utils.ErrDuplicatedVote,
		},
		{
			name:     "Poll Not Found",
			pollID:   "non-existent",
			userID:   "user2",
			optionID: 1,
			wantErr:  utils.ErrPollNotFound,
		},
		{
			name:     "Option Not Found",
			pollID:   pollID,
			userID:   "user3",
			optionID: 99,
			wantErr:  utils.ErrOptionsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{polls: make(map[string]model.Poll)}
			repo.polls[pollID] = model.Poll{
				ID:    pollID,
				Votes: make(map[string]struct{}),
				Options: []model.Option{
					{ID: 1, Text: "Opt 1"},
					{ID: 2, Text: "Opt 2"},
				},
				CreatedAt: time.Now(),
				Duration:  time.Hour,
			}
			if tt.setup != nil {
				tt.setup(repo)
			}
			s := NewService(repo)

			ctx := context.WithValue(context.Background(), middleware.UserInfo("userID"), tt.userID)
			err := s.RegisterVote(ctx, tt.pollID, tt.optionID)
			if err != tt.wantErr {
				t.Errorf("RegisterVote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterVote_Expired(t *testing.T) {
	repo := &mockRepository{polls: make(map[string]model.Poll)}
	s := NewService(repo)

	pollID := "expired-poll"
	repo.polls[pollID] = model.Poll{
		ID:        pollID,
		Votes:     make(map[string]struct{}),
		CreatedAt: time.Now().Add(-2 * time.Hour),
		Duration:  time.Hour,
	}

	ctx := context.WithValue(context.Background(), middleware.UserInfo("userID"), "user1")
	err := s.RegisterVote(ctx, pollID, 1)
	if err != utils.ErrPollExpired {
		t.Errorf("expected ErrPollExpired, got %v", err)
	}
}

func TestListPolls(t *testing.T) {
	repo := &mockRepository{polls: make(map[string]model.Poll)}
	s := NewService(repo)

	repo.polls["p1"] = model.Poll{ID: "p1", Text: "Poll 1"}
	repo.polls["p2"] = model.Poll{ID: "p2", Text: "Poll 2"}

	polls := s.ListPolls(context.Background())
	if len(polls) != 2 {
		t.Errorf("expected 2 polls, got %d", len(polls))
	}
}
