package event

import (
	"strconv"

	"github.com/adinfinit/jamvote/user"
)

type TeamRepo interface {
	CreateTeam(id EventID, team *Team) (TeamID, error)
	UpdateTeam(id EventID, team *Team) error
	TeamByID(id EventID, teamid TeamID) (*Team, error)
	Teams(id EventID) ([]*Team, error)
}

type TeamID int64

func (id TeamID) String() string { return strconv.Itoa(int(id)) }

type Team struct {
	EventID EventID `datastore:"-"`
	ID      TeamID  `datastore:"-"`

	Name    string
	Members []Member
	Game    Game `datastore:",noindex"`
}

type Member struct {
	ID   user.UserID // can be zero
	Name string
}

type Game struct {
	Name string
	Link string `datastore:",noindex"`
	Info string `datastore:",noindex"`
}

func (team *Team) HasEditor(user *user.User) bool {
	if user == nil {
		return false
	}
	if user.Admin {
		return true
	}
	return team.HasMember(user)
}

func (team *Team) HasMember(user *user.User) bool {
	if user == nil {
		return false
	}
	for _, m := range team.Members {
		if m.ID == user.ID {
			return true
		}
	}
	return false
}

func (team *Team) MembersWithEmpty() []Member {
	members := append([]Member{}, team.Members...)
	for len(members) < 6 {
		members = append(members, Member{})
	}
	return members
}
