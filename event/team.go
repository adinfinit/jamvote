package event

import (
	"github.com/adinfinit/jamvote/user"
)

type TeamRepo interface {
	CreateTeam(id EventID, team *Team) (TeamID, error)
	TeamByID(id EventID, teamid TeamID) (*Team, error)
	Teams(id EventID) ([]*Team, error)
}

type TeamID int64

type Team struct {
	ID      TeamID
	Name    string
	Members []user.UserID
	Entry   Entry `datastore:",noindex"`
}

type Member struct {
	ID   user.UserID // can be zero
	Name string      `datastore:",noindex"`
}

type Entry struct {
	Name string
	Info string `datastore:",noindex"`
	Link struct {
		Win string `datastore:",noindex"`
		Mac string `datastore:",noindex"`
		Web string `datastore:",noindex"`
	} `datastore:",noindex"`
}
