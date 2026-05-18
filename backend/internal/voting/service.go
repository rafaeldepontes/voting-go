package voting

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/rafaeldepontes/voting-go/internal/poll/model"
)

type Service interface {

	// ListPolls list all availables polls.
	ListPolls(ctx context.Context) []model.PollDto

	// CreatePoll creates the poll, increases the poll counter
	// and managers the in memory map.
	CreatePoll(ctx context.Context, p model.PollReq) (string, error)

	// CancelPoll cancels an existing poll if the user trying
	// the action is actually the "poll owner". Unauthorized
	// if not.
	CancelPoll(ctx context.Context, pollID string) error

	// RegisterVote updates the in memory poll pointer,
	// manage the option control and then broadcast the info.
	RegisterVote(ctx context.Context, pollID string, optionID int) error

	// Subscribe subscribes a new "user" to our in memory map.
	Subscribe(ctx context.Context, pollID string, conn *websocket.Conn) error
}
