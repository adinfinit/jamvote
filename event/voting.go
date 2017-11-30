package event

import (
	"fmt"
	"net/http"
	"strconv"
)

type Range struct{ Min, Max, Step float64 }

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
		context.Redirect(context.Event.Path("results"), http.StatusSeeOther)
		return
	}

	ballot, err := context.Events.UserBallot(context.Event.ID, context.CurrentUser.ID, context.Team.ID)
	if err != nil && err != ErrNotExists {
		context.FlashNow(err.Error())
	}
	if ballot == nil {
		ballot = &Ballot{}
	}

	ballotinfo := &BallotInfo{
		Team:   context.Team,
		Ballot: ballot,
	}

	context.Data["Aspects"] = AspectsInfo
	context.Data["Ballot"] = ballotinfo

	if context.Request.Method == http.MethodPost {
		if err := context.Request.ParseForm(); err != nil {
			context.FlashNow("Parse form: " + err.Error())
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-vote")
			return
		}

		ballot.Voter = context.CurrentUser.ID
		ballot.Team = context.Team.ID

		readAspect := func(target *Aspect, name string) {
			target.Comment = context.Request.FormValue(name + ".Comment")
			scorestr := context.Request.FormValue(name + ".Score")
			if val, err := strconv.ParseFloat(scorestr, 64); err == nil {
				target.Score = val
			} else {
				context.Flash(name + " value had error: " + err.Error())
			}
		}

		readAspect(&ballot.Theme, "Theme")
		readAspect(&ballot.Enjoyment, "Enjoyment")
		readAspect(&ballot.Aesthetics, "Aesthetics")
		readAspect(&ballot.Innovation, "Innovation")
		readAspect(&ballot.Bonus, "Bonus")

		ballot.Aspects.EnsureRange()
		ballot.Aspects.UpdateTotal()
		ballot.Completed = true

		err := context.Events.SubmitBallot(context.Event.ID, ballot)
		if err != nil {
			context.Flash(err.Error())
		}

		context.Redirect(context.Event.Path("teams"), http.StatusSeeOther)
		return
	}

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
