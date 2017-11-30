package event

import (
	"fmt"

	"github.com/adinfinit/jamvote/user"
	"google.golang.org/appengine/datastore"
)

type BallotRepo interface {
	SubmitBallot(eventid EventID, ballot *Ballot) error
	UserBallots(userid user.UserID, eventid EventID) ([]*Ballot, error)
	UserBallot(userid user.UserID, eventid EventID, teamid TeamID) (*Ballot, error)
	AllBallots(eventid EventID) ([]*Ballot, error)
}

type Ballot struct {
	ID        *datastore.Key `datastore:"-"`
	Voter     user.UserID
	Team      TeamID
	Index     int64 `datastore:",noindex"`
	Submitted bool  `datastore:",noindex"`
	Aspects
}

type BallotInfo struct {
	Team   *Team
	Ballot *Ballot
}

type TeamInfo struct {
	Team     *Team
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
