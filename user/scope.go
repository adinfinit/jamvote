package user

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"google.golang.org/appengine"
)

type Scope struct {
	CurrentUser *User
	Request     *http.Request
	Response    http.ResponseWriter
	Data        map[string]interface{}

	Users Repo

	context.Context
}

func (scope *Scope) Redirect(url string, status int) {
	http.Redirect(scope.Response, scope.Request, url, status)
}

func (scope *Scope) Error(text string, status int) {
	http.Error(scope.Response, text, status)
}

func (scope *Scope) Render(name string) {
	t, err := template.ParseGlob("templates/**/*.html")
	if err != nil {
		http.Error(scope.Response, fmt.Sprintf("Template error: %q", err), http.StatusInternalServerError)
		return
	}

	if err := t.ExecuteTemplate(scope.Response, name+".html", scope.Data); err != nil {
		log.Println(err)
	}
}

var ErrNotExists = errors.New("info does not exist")

func (users *Server) CurrentUser(scope *Scope) *User {
	cred := users.Auth.CurrentCredentials(scope.Request)
	if cred == nil {
		return nil
	}

	user, err := scope.Users.ByCredentials(cred)

	if err == ErrNotExists {
		user = &User{Name: cred.Name}
		_, err := scope.Users.Create(cred, user)
		if err != nil {
			return nil
		}
		return user
	}

	return user
}

func (server *Server) Scope(w http.ResponseWriter, r *http.Request) *Scope {
	scope := &Scope{}
	scope.Context = appengine.NewContext(r)
	scope.Request = r
	scope.Response = w
	scope.Users = &Datastore{scope}
	scope.CurrentUser = server.CurrentUser(scope)
	scope.Data = map[string]interface{}{}
	scope.Data["CurrentUser"] = scope.CurrentUser
	return scope
}

func (server *Server) Scoped(fn func(scope *Scope)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(server.Scope(w, r))
	})
}
