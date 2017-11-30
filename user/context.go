package user

import (
	"net/http"

	"github.com/adinfinit/jamvote/site"
)

type Context struct {
	CurrentUser *User
	Users       Repo

	*site.Context
}

func (server *Server) CurrentUser(context *Context) *User {
	cred := server.Auth.CurrentCredentials(context, context.Request)
	if cred == nil {
		return nil
	}

	user, err := context.Users.ByCredentials(cred)
	if err == ErrNotExists {
		user = &User{Name: cred.Name}
		_, err := context.Users.Create(cred, user)
		if err != nil {
			return nil
		}

		user.NewUser = true
		return user
	}

	// override user rights from credentials
	user.Admin = user.Admin || cred.Admin

	return user
}

func (server *Server) Context(w http.ResponseWriter, r *http.Request) *Context {
	context := &Context{}
	context.Context = server.Site.Context(w, r)
	context.Users = &Datastore{context}
	context.CurrentUser = server.CurrentUser(context)
	context.Data["CurrentUser"] = context.CurrentUser
	return context
}

func (server *Server) Handler(fn func(*Context)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(server.Context(w, r))
	})
}
