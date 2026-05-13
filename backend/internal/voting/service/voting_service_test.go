package service

import (
	"testing"

	"github.com/rafaeldepontes/voting-go/internal/poll/model"
)

func TestCreatePoll(t *testing.T) {
	s := NewService(nil)

	tests := []struct {
		name    string
		req     model.PollReq
		wantErr bool
	}{
		{
			name: "Success",
			req: model.PollReq{
				Name:    "Favorite Color",
				Options: []string{"Red", "Blue", "Green"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := s.CreatePoll(t.Context(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePoll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if id == "" && !tt.wantErr {
				t.Error("CreatePoll() returned empty ID")
			}
		})
	}
}

func TestRegisterVote(t *testing.T) {
	s := NewService(nil)
	pollID, _ := s.CreatePoll(t.Context(), model.PollReq{
		Name:    "Test Poll",
		Options: []string{"Opt 1", "Opt 2"},
	})

	tests := []struct {
		name     string
		pollID   string
		optionID int
		wantErr  bool
	}{
		{
			name:     "Valid Vote",
			pollID:   pollID,
			optionID: 1,
			wantErr:  false,
		},
		{
			name:     "Invalid Poll ID",
			pollID:   "999",
			optionID: 1,
			wantErr:  true,
		},
		{
			name:     "Invalid Option ID",
			pollID:   pollID,
			optionID: 99,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.RegisterVote(t.Context(), tt.pollID, tt.optionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterVote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
