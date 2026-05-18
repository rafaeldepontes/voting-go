package poll

import (
	"context"

	"github.com/rafaeldepontes/voting-go/internal/poll/model"
)

type Repository interface {
	ListPolls(ctx context.Context) ([]model.Poll, error)
	Insert(ctx context.Context, p model.Poll) error
	Update(ctx context.Context, p model.Poll) error
	FindPollByID(ctx context.Context, pollID string) (model.Poll, error)
	Remove(ctx context.Context, pollID string) error
}
