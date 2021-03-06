package profile

import (
	"net/http"

	"github.com/adinfinit/jamvote/event"
	"github.com/adinfinit/jamvote/user"
)

// Context is context for a user.
type Context struct {
	Teams event.TeamRepo
	*user.Context
}

// Context creates a Context for the specified request.
func (server *Server) Context(w http.ResponseWriter, r *http.Request) *Context {
	context := &Context{}
	context.Context = server.Users.Context(w, r)
	context.Teams = server.Events.Events(context)
	return context
}

// Handler wraps automatically fn with Context creation.
func (server *Server) Handler(fn func(*Context)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(server.Context(w, r))
	})
}
