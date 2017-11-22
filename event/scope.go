package event

import (
	"net/http"

	"github.com/adinfinit/jamvote/user"
)

type Scope struct {
	*user.Scope

	Event *Event
}

func (event *Server) Scope(w http.ResponseWriter, r *http.Request) *Scope {
	scope := &Scope{}
	scope.Scope = event.Users.Scope(w, r)

	scope.Event = &Event{
		Slug:  event.Slug,
		Title: event.Title,
	}
	scope.Data["Event"] = scope.Event

	return scope
}

func (event *Server) Scoped(fn func(scope *Scope)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(event.Scope(w, r))
	})
}
