package event

import (
	"encoding/gob"
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

var ErrNotExists = errors.New("does not exist")
var ErrExists = errors.New("already exists")

type EventID string

func (id EventID) String() string { return string(id) }

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
	ID    EventID `datastore:"-"`
	Name  string
	Theme string `datastore:",noindex"`
	Info  string `datastore:",noindex"`

	// Voting allow voting
	Voting bool `datastore:",noindex"`
	// Closed for new entries
	Closed bool `datastore:",noindex"`
	// Revealed, results are publicly viewable
	Revealed bool `datastore:",noindex"`

	VotingOpens  time.Time `datastore:",noindex"`
	VotingCloses time.Time `datastore:",noindex"`

	Organizers []user.UserID `datastore:",noindex"`
	Jammers    []user.UserID `datastore:",noindex"`
}

func init() {
	gob.Register(&Event{})
}

func (event *Event) CanVote() bool {
	return event.Voting && !event.Closed
}

func (event *Event) HasJammer(u *user.User) bool {
	if u == nil {
		return false
	}

	return containsUser(event.Jammers, u.ID)
}

func containsUser(userids []user.UserID, userid user.UserID) bool {
	for _, jammer := range userids {
		if jammer == userid {
			return true
		}
	}
	return false
}

func (event *Event) AddRemoveJammers(added, removed []user.UserID) {
	result := []user.UserID{}
	for _, userid := range event.Jammers {
		if !containsUser(removed, userid) {
			result = append(result, userid)
		}
	}
	for _, userid := range added {
		if !containsUser(result, userid) {
			result = append(result, userid)
		}
	}
	event.Jammers = result
}
