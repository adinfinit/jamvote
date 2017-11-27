package event

import (
	"fmt"

	"github.com/adinfinit/jamvote/user"
)

type BallotRepo interface {
	SubmitBallot(eventid EventID, ballot *Ballot) error
	UserBallots(userid user.UserID, eventid EventID) ([]*Ballot, error)
	UserBallot(userid user.UserID, eventid EventID, teamid TeamID) (*Ballot, error)
	AllBallots(eventid EventID) ([]*Ballot, error)
}

type Ballot struct {
	Voter user.UserID
	Team  TeamID
	Aspects
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
	return fmt.Sprintf("%.2f", aspect.Score)
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

/*

Theme
How well does it interpret the theme
1 Not even close
2 Resembling
3 Related
4 Spot on
5 Novel Interpretation

Enjoyment
How does the game generally feel
1 Boring
2 Not playing again
3 Nice
4 Didn't want to stop
5 Will play later.

Aesthetics
How well is the story, art and audio executed
1 None
2 Needs tweaks
3 Nice
4 Really good
5 Exceptional

Innovation
Something novel in the game
1 Seen it a lot
2 Interesting variation
3 Interesting approach
4 Never seen before
5 Exceptional

Bonus
Anything exceptionally special about * 0,5
1 Nothing special
2 Really liked *
3 Really loved *
4 Loved everything
5 <3

*/
