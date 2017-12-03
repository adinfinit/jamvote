package event

import (
	"fmt"
	"net/http"

	"github.com/adinfinit/jamvote/site"
	"github.com/adinfinit/jamvote/user"
)

type Context struct {
	Event  *Event
	Team   *Team
	Events Repo

	*user.Context
}

func (server *Server) Context(w http.ResponseWriter, r *http.Request) *Context {
	context := &Context{}
	context.Context = server.Users.Context(w, r)
	context.Events = &Datastore{context}

	eventid, ok := context.StringParam("eventid")
	if ok && EventID(eventid).Valid() {
		event, err := context.Events.ByID(EventID(eventid))
		if err == nil && event != nil {
			context.Event = event
			context.Data["Event"] = context.Event
		} else {
			context.FlashErrorNow(err.Error())
		}
	}

	if context.Event != nil {
		if !context.Event.Voting {
			if site.IsValidTime(context.Event.VotingOpens) {
				context.Data["VotingOpens"] = site.NewCountdown(context.Event.VotingOpens)
			}
		} else if !context.Event.Closed {
			if site.IsValidTime(context.Event.VotingCloses) {
				context.Data["VotingCloses"] = site.NewCountdown(context.Event.VotingCloses)
			}
		}
	}

	teamid, ok := context.IntParam("teamid")
	if ok && context.Event != nil {
		team, err := context.Events.TeamByID(context.Event.ID, TeamID(teamid))
		if err == nil && team != nil {
			context.Team = team
			context.Data["Team"] = context.Team
			context.Data["CanEditTeam"] = context.Team.HasEditor(context.CurrentUser)
		}
	}

	return context
}

func (server *Server) HandlerMaybe(fn func(*Context)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(server.Context(w, r))
	})
}

func (server *Server) Handler(fn func(*Context)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := server.Context(w, r)
		if context.Event == nil {
			eventid, _ := context.StringParam("eventid")
			context.FlashError(fmt.Sprintf("Event %q does not exist.", eventid))
			context.Redirect("/", http.StatusSeeOther)
			return
		}
		fn(context)
	})
}
