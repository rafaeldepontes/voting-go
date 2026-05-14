package service

import (
	"context"
	"testing"

	"github.com/rafaeldepontes/voting-go/internal/poll/model"
)

type mockRepository struct{}

func (m *mockRepository) Insert(ctx context.Context, poll model.Poll) error {
	return nil
}

func (m *mockRepository) FindPollByID(ctx context.Context, id string) (model.Poll, error) {
	return model.Poll{}, nil
}

func (m *mockRepository) ListPolls(ctx context.Context) ([]model.Poll, error) {
	return []model.Poll{}, nil
}

func (m *mockRepository) Update(ctx context.Context, poll model.Poll) error {
	return nil
}

func TestCreatePoll(t *testing.T) {
	s := NewService(&mockRepository{})

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
			id, err := s.CreatePoll(context.Background(), tt.req)
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
