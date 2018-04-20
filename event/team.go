package event

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/adinfinit/jamvote/internal/natural"
	"github.com/adinfinit/jamvote/user"
)

type TeamRepo interface {
	CreateTeam(id EventID, team *Team) (TeamID, error)
	UpdateTeam(id EventID, team *Team) error
	TeamByID(id EventID, teamid TeamID) (*Team, error)
	Teams(id EventID) ([]*Team, error)
}

const MaxTeamMembers = 6

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
	Info string `datastore:",noindex"`
	Link struct {
		Jam      string `datastore:",noindex"`
		Download string `datastore:",noindex"`
		Facebook string `datastore:",noindex"`
	} `datastore:",noindex"`
}

func (team *Team) Less(other *Team) bool {
	return natural.Less(team.Name, other.Name)
}

func (team *Team) Verify() error {
	if team.Name == "" {
		return errors.New("Team name cannot be empty.")
	}
	if len(team.Members) == 0 {
		return errors.New("Team must have at least one member.")
	}

	if team.Game.Link.Jam != "" {
		u, err := url.Parse(team.Game.Link.Jam)
		if err != nil {
			return errors.New("Invalid Jam link: " + err.Error())
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return errors.New("Invalid Jam link.")
		}
	}

	if team.Game.Link.Download != "" {
		u, err := url.Parse(team.Game.Link.Download)
		if err != nil {
			return errors.New("Invalid Download link: " + err.Error())
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return errors.New("Invalid Download link.")
		}
	}

	if team.Game.Link.Facebook != "" {
		u, err := url.Parse(team.Game.Link.Facebook)
		if err != nil {
			return errors.New("Invalid Facebook link: " + err.Error())
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return errors.New("Invalid Facebook link.")
		}
	}

	return nil
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
	return team.HasMemberID(user.ID)
}

func (team *Team) HasMemberID(userid user.UserID) bool {
	for _, m := range team.Members {
		if m.ID == userid {
			return true
		}
	}
	return false
}

func (team *Team) HasSubmitted() bool {
	if team.Game.Name == "" {
		return false
	}
	if team.Game.Link.Facebook == "" && team.Game.Link.Jam == "" {
		return false
	}
	return true
}

func (team *Team) MembersWithEmpty() []Member {
	members := append([]Member{}, team.Members...)
	for len(members) < 6 {
		members = append(members, Member{})
	}
	return members
}
