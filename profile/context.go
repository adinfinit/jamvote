package profile

import (
	"net/http"

	"github.com/adinfinit/jamvote/event"
	"github.com/adinfinit/jamvote/user"
)

type Context struct {
	Teams event.TeamRepo
	*user.Context
}

func (server *Server) Context(w http.ResponseWriter, r *http.Request) *Context {
	context := &Context{}
	context.Context = server.Users.Context(w, r)
	context.Teams = &event.Datastore{context}
	return context
}

func (server *Server) Handler(fn func(*Context)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(server.Context(w, r))
	})
}
