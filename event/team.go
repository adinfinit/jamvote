package event

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/adinfinit/jamvote/internal/natural"
	"github.com/adinfinit/jamvote/user"
)

// TeamRepo contains team management in an event.
type TeamRepo interface {
	CreateTeam(id EventID, team *Team) (TeamID, error)
	UpdateTeam(id EventID, team *Team) error
	DeleteTeam(id EventID, teamid TeamID) error
	TeamByID(id EventID, teamid TeamID) (*Team, error)
	Teams(id EventID) ([]*Team, error)
	TeamsByUser(id user.UserID) ([]*EventTeam, error)
}

// MaxTeamMembers defines hard limit on team members.
const MaxTeamMembers = 9

// SuggestedTeamMemberes defines how many members are shown by default.
const SuggestedTeamMembers = 6

// TeamID is a unique identifier for a team.
type TeamID int64

// String is the string representation of the team id.
func (id TeamID) String() string { return strconv.Itoa(int(id)) }

// EventTeam contains information about the event and team.
type EventTeam struct {
	Event Event
	Team
}

// Team contains all information about a team.
type Team struct {
	EventID EventID `datastore:"-"`
	ID      TeamID  `datastore:"-"`

	Name    string
	Members []Member
	Game    Game `datastore:",noindex"`
}

// Member is a team member. There may not be a registered user.
type Member struct {
	ID   user.UserID // can be zero
	Name string
}

// Game contains all information about a game.
type Game struct {
	Name string
	Info string `datastore:",noindex"`

	Noncompeting bool `datastore:",noindex"`

	Link struct {
		Jam      string `datastore:",noindex"`
		Download string `datastore:",noindex"`
		Facebook string `datastore:",noindex"`
	} `datastore:",noindex"`
}

// Less compares teams by name.
func (team *Team) Less(other *Team) bool {
	return natural.Less(team.Name, other.Name)
}

// Verify verifies whether team has valid state.
func (team *Team) Verify() error {
	if team.Name == "" {
		return errors.New("team name cannot be empty")
	}
	if len(team.Members) == 0 {
		return errors.New("team must have at least one member")
	}

	if team.Game.Link.Jam != "" {
		u, err := url.Parse(team.Game.Link.Jam)
		if err != nil {
			return errors.New("invalid Jam link: " + err.Error())
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return errors.New("invalid Jam link")
		}
	}

	if team.Game.Link.Download != "" {
		u, err := url.Parse(team.Game.Link.Download)
		if err != nil {
			return errors.New("invalid Download link: " + err.Error())
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return errors.New("invalid Download link")
		}
	}

	if team.Game.Link.Facebook != "" {
		u, err := url.Parse(team.Game.Link.Facebook)
		if err != nil {
			return errors.New("invalid Facebook link: " + err.Error())
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return errors.New("invalid Facebook link")
		}
	}

	return nil
}

// HasEditor checks whether user can edit the team.
func (team *Team) HasEditor(user *user.User) bool {
	if user == nil {
		return false
	}
	if user.Admin {
		return true
	}
	return team.HasMember(user)
}

// HasMember checks whether user is a member of the team.
func (team *Team) HasMember(user *user.User) bool {
	if user == nil {
		return false
	}
	return team.HasMemberID(user.ID)
}

// HasMemberID checks whether user is a member by userid.
func (team *Team) HasMemberID(userid user.UserID) bool {
	for _, m := range team.Members {
		if m.ID == userid {
			return true
		}
	}
	return false
}

// HasSubmitted checks whether team has all information necessary.
func (team *Team) HasSubmitted() bool {
	if team.Game.Name == "" {
		return false
	}

	hasFacebook := team.Game.Link.Facebook != ""
	hasJam := team.Game.Link.Jam != ""
	hasDownload := team.Game.Link.Download != ""

	return hasFacebook || hasJam || hasDownload
}

// IsCompeting returns whether team is part of the prizes.
func (team *Team) IsCompeting() bool {
	return team.HasSubmitted() && !team.Game.Noncompeting
}

// MembersWithEmpty returns slice with additional empty members if needed.
func (team *Team) MembersWithEmpty() []Member {
	members := append([]Member{}, team.Members...)
	for len(members) < SuggestedTeamMembers {
		members = append(members, Member{})
	}
	return members
}

// MembersForEdit returns slice with additional empty members if needed.
func (team *Team) MembersForEdit(isAdmin bool) []Member {
	members := append([]Member{}, team.Members...)

	maxMembers := SuggestedTeamMembers
	if isAdmin {
		maxMembers = MaxTeamMembers
	}
	for len(members) < maxMembers {
		members = append(members, Member{})
	}
	return members
}
