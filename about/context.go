package about

import (
	"net/http"

	"github.com/adinfinit/jamvote/user"
)

// Context is context for a user.
type Context struct {
	*user.Context
}

// Context creates a Context for the specified request.
func (server *Server) Context(w http.ResponseWriter, r *http.Request) *Context {
	context := &Context{}
	context.Context = server.Users.Context(w, r)
	return context
}

// Handler wraps automatically fn with Context creation.
func (server *Server) Handler(fn func(*Context)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(server.Context(w, r))
	})
}
