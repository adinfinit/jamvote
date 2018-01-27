package event

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
)

type Range struct{ Min, Max, Step float64 }

func (server *Server) FillQueue(context *Context) {
	if context.CurrentUser == nil {
		// TODO: add return address to team-creation page
		context.FlashMessage("You must be logged in to vote.")
		context.Redirect("/user/login", http.StatusSeeOther)
		return
	}

	if !context.Event.HasJammer(context.CurrentUser) {
		context.FlashMessage("You have not been approved for this event.")
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}

	if !context.Event.Voting {
		context.FlashMessage("Voting has not yet started.")
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}

	if context.Event.Closed {
		context.FlashMessage("Voting is closed.")
		if context.Event.Revealed {
			context.Redirect(context.Event.Path("results"), http.StatusSeeOther)
		} else {
			context.Redirect(context.Event.Path(), http.StatusSeeOther)
		}
		return
	}

	_, incomplete, err := context.Events.CreateIncompleteBallots(context.Event.ID, context.CurrentUser.ID)
	if err != nil {
		context.FlashError(err.Error())
	}
	if len(incomplete) == 0 {
		context.FlashMessage("No more games available for now.")
	}

	context.Redirect(context.Event.Path("voting"), http.StatusSeeOther)
}

func (server *Server) Voting(context *Context) {
	if context.CurrentUser == nil {
		// TODO: add return address to team-creation page
		context.FlashMessage("You must be logged in to vote.")
		context.Redirect("/user/login", http.StatusSeeOther)
		return
	}

	if !context.Event.HasJammer(context.CurrentUser) {
		context.FlashMessage("You have not been approved for this event.")
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}

	if !context.Event.Voting {
		context.FlashMessage("Voting has not yet started.")
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}

	ballots, err := context.Events.UserBallots(context.Event.ID, context.CurrentUser.ID)
	if err != nil {
		context.FlashErrorNow(err.Error())
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
	if context.Event.Closed {
		queue = nil
	}

	sort.Slice(completed, func(i, k int) bool {
		return completed[i].Overall.Score > completed[k].Overall.Score
	})

	context.Data["Queue"] = queue
	context.Data["Completed"] = completed

	context.Render("event-voting")
}

func (server *Server) Vote(context *Context) {
	if context.CurrentUser == nil {
		// TODO: add return address to team-creation page
		context.FlashMessage("You must be logged in to vote.")
		context.Redirect("/user/login", http.StatusSeeOther)
		return
	}

	if !context.Event.HasJammer(context.CurrentUser) {
		context.FlashMessage("You have not been approved for this event.")
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}

	if context.Team == nil {
		teamid, _ := context.IntParam("teamid")
		context.FlashMessage(fmt.Sprintf("Team %v does not exist.", teamid))
		context.Redirect(context.Event.Path("voting"), http.StatusSeeOther)
		return
	}

	if !context.Event.Voting {
		context.FlashMessage("Voting has not yet started.")
		context.Redirect(context.Event.Path("voting"), http.StatusSeeOther)
		return
	}

	//if context.Event.Closed {
	//	context.FlashMessage("Voting is closed.")
	//	if context.Event.Revealed {
	//		context.Redirect(context.Event.Path("results"), http.StatusSeeOther)
	//	} else {
	//		context.Redirect(context.Event.Path(), http.StatusSeeOther)
	//	}
	//	return
	//}

	ballot, err := context.Events.UserBallot(context.Event.ID, context.CurrentUser.ID, context.Team.ID)
	if err != nil && err != ErrNotExists {
		context.FlashErrorNow(err.Error())
	}
	if ballot == nil {
		ballot = &Ballot{}
	}

	ballotinfo := &BallotInfo{
		Team:   context.Team,
		Ballot: ballot,
	}

	context.Data["Aspects"] = AspectDescriptions
	context.Data["Ballot"] = ballotinfo

	if context.Request.Method == http.MethodPost {
		if context.Event.Closed {
			context.FlashErrorNow("Voting is closed.")
			context.Response.WriteHeader(http.StatusForbidden)
			context.Render("event-vote")
			return
		}
		if err := context.Request.ParseForm(); err != nil {
			context.FlashErrorNow("Parse form: " + err.Error())
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-vote")
			return
		}

		ballot.Voter = context.CurrentUser.ID
		ballot.Team = context.Team.ID

		readAspect := func(target *Aspect, name string) {
			target.Comment = context.FormValue(name + ".Comment")
			scorestr := context.FormValue(name + ".Score")
			if val, err := strconv.ParseFloat(scorestr, 64); err == nil {
				target.Score = val
			} else {
				context.FlashError(name + " value had error: " + err.Error())
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
			context.FlashError(err.Error())
		}

		context.Redirect(context.Event.Path("voting"), http.StatusSeeOther)
		return
	}

	context.Render("event-vote")
}

func (server *Server) Reveal(context *Context) {
	if !context.Event.Voting {
		context.FlashMessage("Voting has not yet started.")
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}
	if !context.Event.Closed && !context.Event.Revealed {
		context.FlashMessage("Voting is not closed.")
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}
	if !context.CurrentUser.IsAdmin() {
		context.FlashMessage("Only admin can use reveal.")
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}

	results, err := context.Events.Results(context.Event.ID)
	if err != nil {
		context.FlashErrorNow(err.Error())
	}

	sort.Slice(results, func(i, k int) bool {
		return results[i].Average.Overall.Score > results[k].Average.Overall.Score
	})

	if len(results) > 5 {
		results = results[:5]
	}

	context.Data["FullWidth"] = true

	context.Data["Results"] = results
	context.Render("event-reveal")
}

func (server *Server) Results(context *Context) {
	if !context.Event.Voting {
		context.FlashMessage("Voting has not yet started.")
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}

	if !context.Event.Revealed {
		context.FlashMessage("Voting results have not been yet revealed.")
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}

	results, err := context.Events.Results(context.Event.ID)
	if err != nil {
		context.FlashErrorNow(err.Error())
	}

	sort.Slice(results, func(i, k int) bool {
		return results[i].Average.Overall.Score > results[k].Average.Overall.Score
	})

	context.Data["Results"] = results
	context.Render("event-results")
}

func (server *Server) Progress(context *Context) {
	results, err := context.Events.Results(context.Event.ID)
	if err != nil {
		context.FlashErrorNow(err.Error())
	}

	if context.CurrentUser.IsAdmin() {
		context.Data["AutoRefresh"] = 60
	}

	sort.Slice(results, func(i, k int) bool {
		if results[i].HasSubmitted() != results[k].Team.HasSubmitted() {
			return results[i].HasSubmitted()
		}
		return results[i].Team.Name < results[k].Team.Name
	})

	target := 10
	averagePending := 0.0
	averageComplete := 0.0

	max := target * 3 / 2
	for _, result := range results {
		if max < result.Pending {
			max = result.Pending
		}

		averagePending += clamped(float64(result.Pending), 0, float64(target))
		averageComplete += clamped(float64(result.Complete), 0, float64(target))
	}

	context.Data["AveragePending"] = averagePending / float64(len(results))
	context.Data["AverageComplete"] = averageComplete / float64(len(results))

	context.Data["VoteTarget"] = target
	context.Data["VoteMax"] = max
	context.Data["Progress"] = results
	context.Render("event-progress")
}
