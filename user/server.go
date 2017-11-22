package user

import (
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/adinfinit/jamvote/auth"
	"github.com/gorilla/mux"

	"google.golang.org/appengine/datastore"
)

type Renderer interface {
	Render(scope *Scope)
}

type Server struct {
	Auth *auth.Service
}

type Repo interface {
	ByCredentials(cred *auth.Credentials) (*User, error)
	ByID(id ID) (*User, error)

	Create(cred *auth.Credentials, user *User) (ID, error)
	Update(user *User) error
}

func (users *Server) Register(router *mux.Router) {
	router.HandleFunc("/user", users.Scoped(users.Redirect))
	router.HandleFunc("/user/{userid}/edit", users.Scoped(users.Edit))
	router.HandleFunc("/user/login", users.Scoped(users.Login))
	router.HandleFunc("/user/logout", users.Scoped(users.Logout))
	router.HandleFunc("/user/{userid}", users.Scoped(users.Profile))
}

type CredentialsUser struct {
	UserKey *datastore.Key

	Provider string `datastore:",noindex"`
	Email    string `datastore:",noindex"`
	Name     string `datastore:",noindex"`
}

func getUserID(r *http.Request) (ID, bool) {
	s := mux.Vars(r)["userid"]
	if s == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return ID(id), true
}

func (users *Server) Redirect(scope *Scope) {
	if scope.CurrentUser == nil {
		scope.Redirect("/user/login", http.StatusTemporaryRedirect)
		return
	}

	userurl := path.Join("/user", scope.CurrentUser.ID.String(), "edit")
	scope.Redirect(userurl, http.StatusTemporaryRedirect)
}

func (users *Server) Edit(scope *Scope) {
	userid, ok := getUserID(scope.Request)
	if !ok {
		scope.Error("User ID not specified", http.StatusBadRequest)
		return
	}

	user, err := scope.Users.ByID(userid)
	if err != nil {
		scope.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	if !scope.CurrentUser.Equals(user) && !scope.CurrentUser.Admin {
		// access denied
		scope.Redirect(path.Join("/user", user.ID.String()), http.StatusTemporaryRedirect)
		return
	}

	if scope.Request.Method == http.MethodPost {
		if err := scope.Request.ParseForm(); err != nil {
			scope.Error("Parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		name := scope.Request.FormValue("name")
		facebook := scope.Request.FormValue("facebook")
		github := scope.Request.FormValue("github")

		if name != user.Name ||
			facebook != user.Facebook ||
			github != user.Github {

			user.Name = name
			user.Facebook = facebook
			user.Github = github

			err := scope.Users.Update(user)
			if err != nil {
				log.Printf("user.Edit: update %q: %v", userid, err)
				scope.Error("", http.StatusInternalServerError)
				return
			}
		}
	}

	scope.Data["User"] = user
	scope.Render("user-edit")
}

func (users *Server) Profile(scope *Scope) {
	userid, ok := getUserID(scope.Request)
	if !ok {
		scope.Error("User ID not specified", http.StatusBadRequest)
		return
	}

	user, err := scope.Users.ByID(userid)
	if err != nil {
		scope.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	scope.Data["User"] = user
	scope.Render("user-view")
}

func (users *Server) Login(scope *Scope) {
	scope.Data["Logins"] = users.Auth.Links(scope.Request)
	scope.Render("user-login")
}

func (users *Server) Logout(scope *Scope) {
	users.Auth.Logout(scope.Response, scope.Request)
}
