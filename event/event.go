package event

import (
	"time"

	"github.com/adinfinit/jamvote/user"
)

type Stage string

const (
	Draft   Stage = "draft"
	Started       = "started"
	Voting        = "voting"
	Closed        = "closed"
)

type EventID string

type Event struct {
	Slug  EventID `datastore:"-"`
	Title string
	Info  string `datastore:",noindex"`

	Stage   Stage
	Created time.Time
	Started time.Time
	Closed  time.Time

	Organizers []user.UserID
	Teams      []TeamID
}
