package event

import (
	"time"

	"github.com/adinfinit/jamvote/user"
)

type TeamID string

type Member struct {
	ID   user.ID
	Name string
}

type Event struct {
	Name string

	Create time.Time
	Start  time.Time
	Vote   time.Time
	Closed time.Time

	Organizers []user.ID
	Judges     []user.ID
	Teams      []TeamID
}

type Team struct {
	ID      TeamID
	Name    string
	Members []user.ID
	Entry   Entry
}

type Entry struct {
	Name         string
	Instructions string

	Link struct {
		Win string
		Mac string
		Web string
	}
}

type Vote struct {
	ID   user.ID
	Team TeamID

	Aspects  Aspects
	Override bool
	Total    float64
}
