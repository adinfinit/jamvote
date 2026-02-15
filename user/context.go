package user

import (
	"net/http"

	"github.com/adinfinit/jamvote/site"
)

// Context contains the request data for a user page.
type Context struct {
	CurrentUser *User
	Users       Repo

	*site.Context
}

// CurrentUser returns currently logged in user.
func (server *Server) CurrentUser(context *Context) *User {
	cred := server.Auth.CurrentCredentials(context.Request)
	if cred == nil {
		return nil
	}

	user, err := context.Users.ByCredentials(cred)
	if err == ErrNotExists {
		// Try email-based migration for OAuth2 users.
		if cred.Provider == "google" && cred.Email != "" {
			existingID, findErr := context.Users.FindCredentialByEmail(cred.Email)
			if findErr == nil {
				if aliasErr := context.Users.CreateCredentialAlias(cred, existingID); aliasErr == nil {
					if migrated, byErr := context.Users.ByCredentials(cred); byErr == nil {
						migrated.Admin = migrated.Admin || cred.Admin
						return migrated
					}
				}
			}
		}

		user = &User{
			Name:  cred.Name,
			Email: cred.Email,
		}
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

// Context constructs a new user context for a given request.
func (server *Server) Context(w http.ResponseWriter, r *http.Request) *Context {
	context := &Context{}
	context.Context = server.Site.Context(w, r)
	context.Users = server.DB.Users(context)
	context.CurrentUser = server.CurrentUser(context)
	context.Data["CurrentUser"] = context.CurrentUser
	return context
}

// Handler wraps fn with automatic Context creation.
func (server *Server) Handler(fn func(*Context)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(server.Context(w, r))
	})
}
