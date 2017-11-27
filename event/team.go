package event

import (
	"strconv"

	"github.com/adinfinit/jamvote/user"
)

type TeamRepo interface {
	CreateTeam(id EventID, team *Team) (TeamID, error)
	TeamByID(id EventID, teamid TeamID) (*Team, error)
	Teams(id EventID) ([]*Team, error)
}

type TeamID int64

func (id TeamID) String() string { return strconv.Itoa(int(id)) }

type Team struct {
	EventID EventID `datastore:",noindex"`
	ID      TeamID  `datastore:",noindex"`

	Name    string
	Members []Member
	Entry   Entry `datastore:",noindex"`
}

func (team *Team) HasEditor(user *user.User) bool {
	if user.Admin {
		return true
	}
	return team.HasMember(user)
}

func (team *Team) HasMember(user *user.User) bool {
	for _, m := range team.Members {
		if m.ID == user.ID {
			return true
		}
	}
	return false
}

type Member struct {
	ID   user.UserID // can be zero
	Name string
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
