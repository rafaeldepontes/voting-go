package voting

import (
	"github.com/gorilla/websocket"
	"github.com/rafaeldepontes/voting-go/internal/poll/model"
)

type Service interface {

	// ListPolls list all availables polls.
	ListPolls() []model.PollDto

	// CreatePoll creates the poll, increases the poll counter
	// and managers the in memory map.
	CreatePoll(p model.PollReq) (string, error)

	// RegisterVote updates the in memory poll pointer,
	// manage the option control and then broadcast the info.
	RegisterVote(pollID string, optionID int) error

	// Subscribe subscribes a new "user" to our in memory map.
	Subscribe(pollID string, conn *websocket.Conn) error
}
