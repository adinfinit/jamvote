package event

import (
	"fmt"

	"google.golang.org/appengine/v2/datastore"

	"github.com/adinfinit/jamvote/user"
)

// BallotRepo is used to manage ballots for an event.
type BallotRepo interface {
	Ballots(eventid EventID) ([]*Ballot, error)
	CreateIncompleteBallots(eventid EventID, userid user.UserID) (complete, incomplete []*BallotInfo, err error)
	SubmitBallot(eventid EventID, ballot *Ballot) error
	UserBallot(eventid EventID, userid user.UserID, teamid TeamID) (*Ballot, error)
	UserBallots(eventid EventID, userid user.UserID) ([]*BallotInfo, error)
	Results(eventid EventID) ([]*TeamResult, error)
	TeamBallots(eventid EventID, teamid TeamID) ([]*Ballot, error)
}

// Ballot is all information for a single ballot.
type Ballot struct {
	ID        *datastore.Key `datastore:"-"`
	Voter     user.UserID
	Team      TeamID
	Index     int64 `datastore:",noindex"`
	Completed bool  `datastore:",noindex"`
	Aspects
}

// BallotInfo is a single ballot, but contains a reference to the target team.
type BallotInfo struct {
	*Team
	*Ballot
}

// TeamResult contains all information about a single teams result.
type TeamResult struct {
	*Team
	Ballots []*Ballot

	Average       Aspects
	JudgeAverage  Aspects
	JammerAverage Aspects

	Pending  int
	Complete int

	MemberBallots []*Ballot
}

// HasReviewer checks whether team results contains userid.
func (info *TeamResult) HasReviewer(userid user.UserID) bool {
	for _, ballot := range info.Ballots {
		if ballot.Voter == userid {
			return true
		}
	}
	return false
}

// DefaultAspects contains defaults for aspects.
var DefaultAspects = Aspects{
	Theme:      Aspect{3, ""},
	Enjoyment:  Aspect{3, ""},
	Aesthetics: Aspect{3, ""},
	Innovation: Aspect{3, ""},
	Bonus:      Aspect{0, ""},
	Overall:    Aspect{0, ""},
}

// AverageScores returns averages for all aspects.
func AverageScores(ballots []*Ballot, event *Event) (final, jammers, judges Aspects) {
	judgeCount := 0.0
	jammerCount := 0.0
	for _, ballot := range ballots {
		if !ballot.Completed {
			continue
		}

		if event.JudgePercentage > 0 && event.HasJudgeById(&ballot.Voter) {
			judges.Add(&ballot.Aspects)
			judgeCount += 1.0
		} else {
			jammers.Add(&ballot.Aspects)
			jammerCount += 1.0
		}
	}

	if jammerCount > 0 {
		jammers.Scale(1 / jammerCount)
	}
	if judgeCount > 0 {
		judges.Scale(1 / judgeCount)
	}

	if event.JudgePercentage == 0 {
		final = jammers
		return
	}

	p := event.JudgePercentage / 100
	jammerPart := jammers
	jammerPart.Scale(1 - p)
	final.Add(&jammerPart)

	judgesPart := judges
	judgesPart.Scale(p)
	final.Add(&judgesPart)

	return final, jammers, judges
}

// Aspects contains criteria for scoring a game.
type Aspects struct {
	Theme      Aspect
	Enjoyment  Aspect
	Aesthetics Aspect
	Innovation Aspect
	Bonus      Aspect
	Overall    Aspect
}

// Aspect is a single criteria with an optional comment.
type Aspect struct {
	Score   float64
	Comment string
}

// AspectsInfo contains all scores for aspects.
type AspectsInfo struct {
	Theme      AspectInfo
	Enjoyment  AspectInfo
	Aesthetics AspectInfo
	Innovation AspectInfo
	Bonus      AspectInfo
	Overall    AspectInfo
}

// AspectInfo contains all scores for an aspect.
type AspectInfo struct {
	Scores       []float64
	MemberScores []float64
	Comments     []string
}

// String pretty prints an aspect.
func (aspect Aspect) String() string {
	return fmt.Sprintf("%.1f", aspect.Score)
}

// ClearComments clears all comments in aspects.
func (aspects *Aspects) ClearComments() {
	aspects.Theme.Comment = ""
	aspects.Enjoyment.Comment = ""
	aspects.Aesthetics.Comment = ""
	aspects.Innovation.Comment = ""
	aspects.Bonus.Comment = ""
	aspects.Overall.Comment = ""
}

// Item fetches an aspect by name.
func (aspects *AspectsInfo) Item(name string) AspectInfo {
	switch name {
	case "Theme":
		return aspects.Theme
	case "Enjoyment":
		return aspects.Enjoyment
	case "Aesthetics":
		return aspects.Aesthetics
	case "Innovation":
		return aspects.Innovation
	case "Bonus":
		return aspects.Bonus
	case "Overall":
		return aspects.Overall
	}
	return AspectInfo{}
}

// Add adds together two aspects.
func (aspects *Aspects) Add(other *Aspects) {
	aspects.Theme.Score += other.Theme.Score
	aspects.Enjoyment.Score += other.Enjoyment.Score
	aspects.Aesthetics.Score += other.Aesthetics.Score
	aspects.Innovation.Score += other.Innovation.Score
	aspects.Bonus.Score += other.Bonus.Score
	aspects.Overall.Score += other.Overall.Score
}

func (aspects *Aspects) Scale(multiplier float64) {
	aspects.Theme.Score *= multiplier
	aspects.Enjoyment.Score *= multiplier
	aspects.Aesthetics.Score *= multiplier
	aspects.Innovation.Score *= multiplier
	aspects.Bonus.Score *= multiplier
	aspects.Overall.Score *= multiplier
}

// Add includes other into aspects.
func (aspects *AspectsInfo) Add(other *Aspects, isMember bool) {
	aspects.Theme.Add(&other.Theme, isMember)
	aspects.Enjoyment.Add(&other.Enjoyment, isMember)
	aspects.Aesthetics.Add(&other.Aesthetics, isMember)
	aspects.Innovation.Add(&other.Innovation, isMember)
	aspects.Bonus.Add(&other.Bonus, isMember)
	aspects.Overall.Add(&other.Overall, isMember)
}

// Add includes other into aspect.
func (aspect *AspectInfo) Add(other *Aspect, isMember bool) {
	if isMember {
		aspect.MemberScores = append(aspect.MemberScores, other.Score)
	} else {
		aspect.Scores = append(aspect.Scores, other.Score)
	}
	if other.Comment != "" {
		aspect.Comments = append(aspect.Comments, other.Comment)
	}
}

// EnsureRange ensures that all scores are in the appropriate ranges.
func (aspects *Aspects) EnsureRange() {
	clamp(&aspects.Theme.Score, 1, 5)
	clamp(&aspects.Enjoyment.Score, 1, 5)
	clamp(&aspects.Aesthetics.Score, 1, 5)
	clamp(&aspects.Innovation.Score, 1, 5)
	clamp(&aspects.Bonus.Score, 0, 2.5)
	clamp(&aspects.Overall.Score, 0, 5)
}

// Score returns an aspect score based on a name.
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

// Comment returns an aspect comment based on a name.
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

// UpdateTotal updates overall score.
func (aspects *Aspects) UpdateTotal() {
	aspects.Overall.Score = aspects.Total()
}

// Total calculates total score.
func (aspects *Aspects) Total() float64 {
	return clamped((aspects.Theme.Score+
		aspects.Enjoyment.Score+
		aspects.Aesthetics.Score+
		aspects.Innovation.Score+
		aspects.Bonus.Score)/4.5, 1, 5)
}

// AspectNames contains all names of aspects.
var AspectNames = []string{
	"Theme",
	"Enjoyment",
	"Aesthetics",
	"Innovation",
	"Bonus",
	"Overall",
}

// AspectDescription describes an aspect.
type AspectDescription struct {
	Name        string
	Description string
	Range
	Options []string
}

// AspectDescriptions contains information about aspects.
var AspectDescriptions = []AspectDescription{
	{
		Name:        "Theme",
		Description: "How well does it interpret the theme?",
		Range:       Range{Min: 1, Max: 5, Step: 0.1},
		Options:     []string{"Not even close", "Resembling", "Related", "Spot on", "Novel Interpretation"},
	}, {
		Name:        "Enjoyment",
		Description: "How does the game generally feel?",
		Range:       Range{Min: 1, Max: 5, Step: 0.1},
		Options:     []string{"I want my time back", "Boring", "Nice", "Didn't want to stop", "Will play later"},
	}, {
		Name:        "Aesthetics",
		Description: "How well is the story, art and audio executed?",
		Range:       Range{Min: 1, Max: 5, Step: 0.1},
		Options:     []string{"None", "Needs tweaks", "Nice", "Really good", "Exceptional"},
	}, {
		Name:        "Innovation",
		Description: "Something novel in the game?",
		Range:       Range{Min: 1, Max: 5, Step: 0.1},
		Options:     []string{"Seen it a lot", "Interesting variation", "Interesting approach", "Never seen this before", "Exceptional"},
	}, {
		Name:        "Bonus",
		Description: "Anything exceptionally special about it?",
		Range:       Range{Min: 0, Max: 2.5, Step: 0.1},
		Options:     []string{"Nothing special", "Really liked *", "Really loved **"},
	},
}

// AspectDescriptionsWithOverall also includes the overall.
var AspectDescriptionsWithOverall = append(AspectDescriptions,
	AspectDescription{
		Name:        "Overall",
		Description: "Weighted average of topics.",
		Range:       Range{Min: 1, Max: 5, Step: 0.1},
	})

// Clamp forces v to be between min and max.
func clamp(v *float64, min, max float64) {
	if *v < min {
		*v = min
	}
	if *v > max {
		*v = max
	}
}

// clamped returns v such that it is between min and max.
func clamped(v float64, min, max float64) float64 {
	if v < min {
		return min
	} else if v > max {
		return max
	}
	return v
}
