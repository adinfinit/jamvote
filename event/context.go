package event

import (
	"net/http"
	"time"

	"github.com/adinfinit/jamvote/user"
)

type Context struct {
	Event  *Event
	Team   *Team
	Events Repo

	*user.Context
}

var EventStub = Event{
	Name: "Ludum Dare 40",
	Info: "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Dolore nisi, iusto amet dignissimos expedita alias libero temporibus earum? Mollitia dolor illum atque commodi voluptas dicta nulla maiores dolores aspernatur. Quasi!",

	Stage:   Voting,
	Created: time.Now(),
	Started: time.Now(),
	Closed:  time.Now(),
}

var TeamStub = Team{
	ID:   123,
	Name: "Shroomy",
	Members: []Member{
		{0, "Magic Mike"},
		{5275456790069248, "Wolfram"},
	},
}

func (event *Server) Context(w http.ResponseWriter, r *http.Request) *Context {
	context := &Context{}
	context.Context = event.Users.Context(w, r)
	context.Events = &Datastore{context}

	eventid, ok := context.StringParam("eventid")
	if ok && EventID(eventid).Valid() {
		event, err := context.Events.ByID(EventID(eventid))
		if err == nil && event != nil {
			context.Event = event
			context.Data["Event"] = context.Event
		} else {
			stub := EventStub
			stub.ID = EventID(eventid)
			context.Event = &stub
			context.Data["Event"] = context.Event
		}
	}

	teamid, ok := context.IntParam("teamid")
	if ok && context.Event != nil {
		team, err := context.Events.TeamByID(context.Event.ID, TeamID(teamid))
		if err == nil && team != nil {
			context.Team = team
			context.Data["Team"] = context.Team
			context.Data["CanEditTeam"] = context.Team.HasEditor(context.CurrentUser)
		} else {
			stub := TeamStub
			stub.ID = TeamID(teamid)
			context.Team = &stub
			context.Data["Team"] = context.Team
			context.Data["CanEditTeam"] = context.Team.HasEditor(context.CurrentUser)
		}
	}

	return context
}

func (event *Server) HandlerMaybe(fn func(*Context)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(event.Context(w, r))
	})
}

func (event *Server) Handler(fn func(*Context)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := event.Context(w, r)
		if context.Event == nil {
			// TODO: flash invalid event
			context.Redirect("/", http.StatusTemporaryRedirect)
			return
		}
		fn(context)
	})
}
