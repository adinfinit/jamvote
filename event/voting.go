package event

import (
	"fmt"
	"net/http"
)

var AspectsInfo = []struct {
	Name        string
	Description string
	Options     []string
}{
	{
		Name:        "Theme",
		Description: "How well does it interpret the theme?",
		Options:     []string{"Not even close", "Resembling", "Related", "Spot on", "Novel Interpretation"},
	}, {
		Name:        "Enjoyment",
		Description: "How does the game generally feel?",
		Options:     []string{"Boring", "Not playing again", "Nice", "Didn't want to stop", "Will play later"},
	}, {
		Name:        "Aesthetics",
		Description: "How well is the story, art and audio executed?",
		Options:     []string{"None", "Needs tweaks", "Nice", "Really good", "Exceptional"},
	}, {
		Name:        "Innovation",
		Description: "Something novel in the game?",
		Options:     []string{"Seen it a lot", "Interesting variation", "Interesting approach", "Never seen before", "Exceptional"},
	}, {
		Name:        "Bonus",
		Description: "Anything exceptionally special about it?",
		Options:     []string{"Nothing special", "Really liked *", "Really loved *", "Loved everything", "<3"},
	},
}

func (event *Server) StartVoting(context *Context) {
	if context.CurrentUser == nil {
		// TODO: add return address to team-creation page
		context.Flash("You must be logged in to vote.")
		context.Redirect("/user/login", http.StatusSeeOther)
		return
	}

	if !context.Event.Voting {
		context.Flash("Voting has not yet started.")
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}

	if context.Event.Closed {
		context.FlashNow("Voting has been closed.")
		context.Redirect(context.Event.Path("results"), http.StatusSeeOther)
		return
	}

	_, _, err := context.Events.CreateIncompleteBallots(context.Event.ID, context.CurrentUser.ID)
	if err != nil {
		context.Flash(err.Error())
	}

	context.Redirect(context.Event.Path("voting"), http.StatusSeeOther)
}

func (event *Server) Voting(context *Context) {
	if context.CurrentUser == nil {
		// TODO: add return address to team-creation page
		context.Flash("You must be logged in to vote.")
		context.Redirect("/user/login", http.StatusSeeOther)
		return
	}

	if !context.Event.Voting {
		context.Flash("Voting has not yet started.")
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}

	if context.Event.Closed {
		context.FlashNow("Voting has been closed.")
		context.Redirect(context.Event.Path("results"), http.StatusSeeOther)
		return
	}

	ballots, err := context.Events.UserBallots(context.Event.ID, context.CurrentUser.ID)
	if err != nil {
		context.FlashNow(err.Error())
	}

	queue := []*BallotInfo{}
	completed := []*BallotInfo{}

	for _, ballot := range ballots {
		if ballot.Completed {
			completed = append(completed, ballot)
		} else {
			queue = append(queue, ballot)
		}
	}
	if len(ballots) <= 3 {
		queue = ballots
	}

	context.Data["Queue"] = queue
	context.Data["Completed"] = completed

	context.Render("event-voting")
}

func (event *Server) Vote(context *Context) {
	if context.CurrentUser == nil {
		// TODO: add return address to team-creation page
		context.Flash("You must be logged in to vote.")
		context.Redirect("/user/login", http.StatusSeeOther)
		return
	}

	if context.Team == nil {
		teamid, _ := context.IntParam("teamid")
		context.Flash(fmt.Sprintf("Team %v does not exist.", teamid))
		context.Redirect(context.Event.Path("voting"), http.StatusSeeOther)
		return
	}

	if !context.Event.Voting {
		context.Flash("Voting has not yet started.")
		context.Redirect(context.Event.Path("voting"), http.StatusSeeOther)
		return
	}

	if context.Event.Closed {
		context.Flash("Voting has been closed.")
		context.Redirect(context.Event.Path("voting"), http.StatusSeeOther)
		return
	}

	ballot, err := context.Events.UserBallot(context.Event.ID, context.CurrentUser.ID, context.Team.ID)
	if err == ErrNotExists {
		ballot = &Ballot{}
	}

	ballotinfo := &BallotInfo{
		Team:   context.Team,
		Ballot: ballot,
	}

	context.Data["Aspects"] = AspectsInfo
	context.Data["Ballot"] = ballotinfo

	context.Render("event-vote")
}

func (event *Server) Results(context *Context) {
	if !context.CurrentUser.IsAdmin() {
		if !context.Event.Revealed {
			context.Flash("Results have not yet been revealed.")
			context.Redirect(context.Event.Path(), http.StatusSeeOther)
			return
		}
	}

	context.Render("todo")
}
