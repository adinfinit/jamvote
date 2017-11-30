package event

import (
	"errors"
	"fmt"

	"github.com/adinfinit/jamvote/user"
	"google.golang.org/appengine/datastore"
)

type BallotRepo interface {
	CreateIncompleteBallots(eventid EventID, userid user.UserID) (complete, incomplete []*BallotInfo, err error)
	SubmitBallot(eventid EventID, ballot *Ballot) error
	UserBallot(eventid EventID, userid user.UserID, teamid TeamID) (*Ballot, error)
	UserBallots(eventid EventID, userid user.UserID) ([]*BallotInfo, error)
	Results(eventid EventID) ([]*TeamResult, error)
}

var ErrUnfinished = errors.New("Incomplete ballots.")

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

type TeamResult struct {
	*Team
	Ballots  []*Ballot
	Average  Aspects
	Pending  int
	Complete int

	MemberBallots []*Ballot
}

func (info *TeamResult) HasReviewer(userid user.UserID) bool {
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
	Bonus:      Aspect{0, ""},
	Overall:    Aspect{0, ""},
}

func AverageScores(ballots []*Ballot) Aspects {
	count := 0.0
	average := Aspects{}
	for _, ballot := range ballots {
		if !ballot.Completed {
			continue
		}

		if count == 0 {
			average = ballot.Aspects
		} else {
			average.Add(&ballot.Aspects)
		}
		count += 1.0
	}

	if count > 0 {
		average.Theme.Score /= count
		average.Enjoyment.Score /= count
		average.Aesthetics.Score /= count
		average.Innovation.Score /= count
		average.Bonus.Score /= count
		average.Overall.Score /= count
	}

	return average
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

func (aspect Aspect) String() string {
	return fmt.Sprintf("%.1f", aspect.Score)
}

func (aspects *Aspects) ClearComments() {
	aspects.Theme.Comment = ""
	aspects.Enjoyment.Comment = ""
	aspects.Aesthetics.Comment = ""
	aspects.Innovation.Comment = ""
	aspects.Bonus.Comment = ""
	aspects.Overall.Comment = ""
}

func (aspects *Aspects) Add(other *Aspects) {
	aspects.Theme.Score += other.Theme.Score
	aspects.Enjoyment.Score += other.Enjoyment.Score
	aspects.Aesthetics.Score += other.Aesthetics.Score
	aspects.Innovation.Score += other.Innovation.Score
	aspects.Bonus.Score += other.Bonus.Score
	aspects.Overall.Score += other.Overall.Score
}

func (aspects *Aspects) EnsureRange() {
	clamp(&aspects.Theme.Score, 1, 5)
	clamp(&aspects.Enjoyment.Score, 1, 5)
	clamp(&aspects.Aesthetics.Score, 1, 5)
	clamp(&aspects.Innovation.Score, 1, 5)
	clamp(&aspects.Bonus.Score, 0, 2.5)
	clamp(&aspects.Overall.Score, 0, 5)
}

func (aspects *Aspects) Score(name string) float64 {
	switch name {
	case "Theme":
		return aspects.Theme.Score
	case "Enjoyment":
		return aspects.Enjoyment.Score
	case "Aesthetics":
		return aspects.Aesthetics.Score
	case "Innovation":
		return aspects.Innovation.Score
	case "Bonus":
		return aspects.Bonus.Score
	case "Overall":
		return aspects.Overall.Score
	}
	return 0
}

func (aspects *Aspects) Comment(name string) string {
	switch name {
	case "Theme":
		return aspects.Theme.Comment
	case "Enjoyment":
		return aspects.Enjoyment.Comment
	case "Aesthetics":
		return aspects.Aesthetics.Comment
	case "Innovation":
		return aspects.Innovation.Comment
	case "Bonus":
		return aspects.Bonus.Comment
	case "Overall":
		return aspects.Overall.Comment
	}
	return "INVALID"
}

func (aspects *Aspects) UpdateTotal() {
	aspects.Overall.Score = aspects.Total()
}

func (aspects *Aspects) Total() float64 {
	return clamped((aspects.Theme.Score+
		aspects.Enjoyment.Score+
		aspects.Aesthetics.Score+
		aspects.Innovation.Score+
		aspects.Bonus.Score)/4.5, 1, 5)
}

var AspectsInfo = []struct {
	Name        string
	Description string
	Range
	Options []string
}{
	{
		Name:        "Theme",
		Description: "How well does it interpret the theme?",
		Range:       Range{Min: 1, Max: 5, Step: 0.1},
		Options:     []string{"Not even close", "Resembling", "Related", "Spot on", "Novel Interpretation"},
	}, {
		Name:        "Enjoyment",
		Description: "How does the game generally feel?",
		Range:       Range{Min: 1, Max: 5, Step: 0.1},
		Options:     []string{"Boring", "Not playing again", "Nice", "Didn't want to stop", "Will play later"},
	}, {
		Name:        "Aesthetics",
		Description: "How well is the story, art and audio executed?",
		Range:       Range{Min: 1, Max: 5, Step: 0.1},
		Options:     []string{"None", "Needs tweaks", "Nice", "Really good", "Exceptional"},
	}, {
		Name:        "Innovation",
		Description: "Something novel in the game?",
		Range:       Range{Min: 1, Max: 5, Step: 0.1},
		Options:     []string{"Seen it a lot", "Interesting variation", "Interesting approach", "Never seen before", "Exceptional"},
	}, {
		Name:        "Bonus",
		Description: "Anything exceptionally special about it?",
		Range:       Range{Min: 0, Max: 2.5, Step: 0.1},
		Options:     []string{"Nothing special", "Really liked *", "Really loved **"},
	},
}

func clamp(v *float64, min, max float64) {
	if *v < min {
		*v = min
	}
	if *v > max {
		*v = max
	}
}

func clamped(v float64, min, max float64) float64 {
	if v < min {
		return min
	} else if v > max {
		return max
	}
	return v
}
