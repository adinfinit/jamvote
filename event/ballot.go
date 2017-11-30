package event

import (
	"fmt"

	"github.com/adinfinit/jamvote/user"
	"google.golang.org/appengine/datastore"
)

type BallotRepo interface {
	CreateIncompleteBallots(eventid EventID, userid user.UserID) (complete, incomplete []*BallotInfo, err error)
	SubmitBallot(eventid EventID, ballot *Ballot) error
	UserBallot(eventid EventID, userid user.UserID, teamid TeamID) (*Ballot, error)
	UserBallots(eventid EventID, userid user.UserID) ([]*BallotInfo, error)
	AllTeamInfos(eventid EventID) ([]*TeamInfo, error)
}

type Ballot struct {
	ID        *datastore.Key `datastore:"-"`
	Voter     user.UserID
	Team      TeamID
	Index     int64 `datastore:",noindex"`
	Completed bool  `datastore:",noindex"`
	Aspects
}

type BallotInfo struct {
	*Team
	*Ballot
}

type TeamInfo struct {
	*Team
	Ballots  []*Ballot
	Pending  int
	Complete int
}

func (info *TeamInfo) HasReviewer(userid user.UserID) bool {
	for _, ballot := range info.Ballots {
		if ballot.Voter == userid {
			return true
		}
	}
	return false
}

var DefaultAspects = Aspects{
	Theme:      Aspect{3, ""},
	Enjoyment:  Aspect{3, ""},
	Aesthetics: Aspect{3, ""},
	Innovation: Aspect{3, ""},
	Bonus:      Aspect{3, ""},
	Overall:    Aspect{3, ""},
}

type Aspects struct {
	Theme      Aspect
	Enjoyment  Aspect
	Aesthetics Aspect
	Innovation Aspect
	Bonus      Aspect
	Overall    Aspect
}

type Aspect struct {
	Score   float64
	Comment string
}

func (aspect *Aspect) String() string {
	return fmt.Sprintf("%.1f", aspect.Score)
}

func (aspects *Aspects) EnsureRange() {
	clamp(&aspects.Theme.Score, 1, 5)
	clamp(&aspects.Enjoyment.Score, 1, 5)
	clamp(&aspects.Aesthetics.Score, 1, 5)
	clamp(&aspects.Innovation.Score, 1, 5)
	clamp(&aspects.Bonus.Score, 1, 5)
}

func (aspects *Aspects) Total() float64 {
	return (aspects.Theme.Score +
		aspects.Enjoyment.Score +
		aspects.Aesthetics.Score +
		aspects.Innovation.Score +
		aspects.Bonus.Score*0.5) / (5*4 + 5*0.5)
}

func clamp(v *float64, min, max float64) {
	if *v < min {
		*v = min
	}
	if *v > max {
		*v = max
	}
}
