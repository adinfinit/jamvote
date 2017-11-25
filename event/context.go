package event

import (
	"net/http"

	"github.com/adinfinit/jamvote/user"
)

type Context struct {
	Event  *Event
	Events Repo

	*user.Context
}

func (event *Server) Context(w http.ResponseWriter, r *http.Request) *Context {
	context := &Context{}
	context.Context = event.Users.Context(w, r)
	context.Events = &Datastore{context}

	id, ok := context.IntParam("eventid")
	if ok && EventID(id).Valid() {
		event, err := context.Events.ByID(EventID(id))
		if err == nil && event != nil {
			context.Event = event
			context.Data["Event"] = context.Event
		}
	}

	return context
}

func (event *Server) Handler(fn func(*Context)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := event.Context(w, r)
		if ctx.Event == nil {
			// TODO: flash invalid event-id
			ctx.Redirect("/", http.StatusTemporaryRedirect)
			return
		}

		fn(ctx)
	})
}
