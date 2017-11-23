package event

import (
	"net/http"

	"github.com/adinfinit/jamvote/user"
)

type Context struct {
	Event *Event

	*user.Context
}

func (event *Server) Context(w http.ResponseWriter, r *http.Request) *Context {
	context := &Context{}
	context.Context = event.Users.Context(w, r)

	context.Event = &Event{
		Slug:  event.Slug,
		Title: event.Title,
	}
	context.Data["Event"] = context.Event

	return context
}

func (event *Server) Handler(fn func(*Context)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(event.Context(w, r))
	})
}
