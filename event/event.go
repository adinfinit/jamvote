package event

import (
	"errors"
	"time"

	"github.com/adinfinit/jamvote/user"
)

type Repo interface {
	List() ([]*Event, error)

	Create(event *Event) error
	ByID(id EventID) (*Event, error)
	Update(event *Event) error

	TeamRepo
	BallotRepo
}

var ErrNotExists = errors.New("info does not exist")

type Stage string

const (
	Draft   Stage = "draft"
	Started       = "started"
	Voting        = "voting"
	Closed        = "closed"
)

type EventID string

func (id EventID) Valid() bool {
	if id == "" {
		return false
	}

	for _, r := range id {
		switch {
		case 'a' <= r && r <= 'z':
		case '0' <= r && r <= '9':
		case '-' == r:
		default:
			return false
		}
	}

	return true
}

type Event struct {
	Slug EventID `datastore:"-"`
	Name string
	Info string `datastore:",noindex"`

	Stage   Stage
	Created time.Time
	Started time.Time
	Closed  time.Time

	Organizers []user.UserID
}
