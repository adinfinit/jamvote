package event

import (
	"context"
	"encoding/gob"
	"errors"
	"time"

	"github.com/adinfinit/jamvote/user"
)

// DB is the master database.
type DB interface {
	Events(context context.Context) Repo
}

// Repo describes interaction with the database.
type Repo interface {
	List() ([]*Event, error)

	Create(event *Event) error
	ByID(id EventID) (*Event, error)
	Update(event *Event) error

	TeamRepo
	BallotRepo
}

// ErrNotExists is returned when an event doesn't exist.
var ErrNotExists = errors.New("does not exist")
// ErrExists is returned when an event already exists.
var ErrExists = errors.New("already exists")

// EventID is a unique identifier for an event.
type EventID string

// String returns printable event id.
func (id EventID) String() string { return string(id) }

// Valid checks whether event is valid.
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

// Event contains all information about an event.
type Event struct {
	ID    EventID `datastore:"-"`
	Name  string
	Theme string `datastore:",noindex"`
	Info  string `datastore:",noindex"`

	// Created is the time when the event was created
	Created time.Time `datastore:",noindex"`
	// StartTime is the starting time of the event
	StartTime time.Time `datastore:",noindex"`
	// EndTime is the end time of the event
	EndTime time.Time `datastore:",noindex"`

	// New Registration is allowed
	Registration bool `datastore:",noindex"`
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

// CanVote returns whether it's possible to vote in this event.
func (event *Event) CanVote() bool {
	return event.Voting && !event.Closed
}

// CanRegister returns whether u can register to the event.
func (event *Event) CanRegister(u *user.User) bool {
	if u.IsAdmin() {
		return true
	}
	return !event.Closed && event.Registration
}

// HasJammer checks whether u has registered.
func (event *Event) HasJammer(u *user.User) bool {
	if u == nil {
		return false
	}

	return containsUser(event.Jammers, u.ID)
}

// containsUser checks whether userids contains userid.
func containsUser(userids []user.UserID, userid user.UserID) bool {
	for _, jammer := range userids {
		if jammer == userid {
			return true
		}
	}
	return false
}

// AddRemoveJammers adds and removes jammers from the event.
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

// Less compares events based on start time.
func (event *Event) Less(other *Event) bool {
	return event.startTime().After(other.startTime())
}

// startTime returns event start time.
func (event *Event) startTime() time.Time {
	if !event.StartTime.IsZero() {
		return event.StartTime
	}
	return event.Created
}
