package user

import (
	"net/http"

	"github.com/adinfinit/jamvote/auth"
	"github.com/adinfinit/jamvote/site"
	"github.com/gorilla/mux"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type Server struct {
	Auth     *auth.Service
	Renderer *site.Renderer
}

func (users *Server) Register(router *mux.Router) {
	router.HandleFunc("/user", users.EditProfile) // REDIRECT TO SELF
	router.HandleFunc("/user/login", users.Login)
	router.HandleFunc("/user/logout", users.Logout)
	router.HandleFunc("/user/{userid}", users.Profile)
}

type CredentialsUser struct {
	UserID ID
}

func (users *Server) EditProfile(w http.ResponseWriter, r *http.Request) {
	cred := users.Auth.CurrentCredentials(w, r)
	if cred == nil {
		http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)
		return
	}

	c := appengine.NewContext(r)

	authkey := datastore.NewKey(c, "Auth", cred.ID, 0, nil)

	creduser := &CredentialsUser{}
	err := datastore.Get(c, authkey, creduser)

	if err == datastore.ErrNoSuchEntity {
		// create new user
		user := &User{
			ID:    ID(cred.ID),
			Name:  cred.Name,
			Email: cred.Email,
		}

		// create new user
		userkey := datastore.NewKey(c, "User", string(user.ID), 0, nil)
		if _, err := datastore.Put(c, userkey, user); err != nil {
			http.Error(w, "Create new user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		creduser.UserID = user.ID

		// add authentication mapping
		if _, err := datastore.Put(c, authkey, creduser); err != nil {
			http.Error(w, "Create new auth: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	user := &User{}
	userkey := datastore.NewKey(c, "User", string(creduser.UserID), 0, nil)

	if err := datastore.Get(c, userkey, user); err != nil {
		http.Error(w, "Get user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	user.ID = creduser.UserID

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		facebook := r.FormValue("facebook")

		if name != user.Name || email != user.Email || facebook != user.Facebook {
			user.Name = name
			user.Email = email
			user.Facebook = facebook

			if _, err := datastore.Put(c, userkey, user); err != nil {
				http.Error(w, "Put user: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	users.Renderer.Render(w, "user-edit", map[string]interface{}{
		"User": user,
	})
}

func (users *Server) Profile(w http.ResponseWriter, r *http.Request) {
	userid := mux.Vars(r)["userid"]
	users.Renderer.Render(w, "user-edit", map[string]interface{}{
		"UserID": userid,
	})
}

func (users *Server) Login(w http.ResponseWriter, r *http.Request) {
	type Login struct{ Title, URL string }

	logins := users.Auth.Logins(r)
	users.Renderer.Render(w, "user-login", map[string]interface{}{
		"Logins": logins,
	})
}

func (users *Server) Logout(w http.ResponseWriter, r *http.Request) {
	users.Auth.Logout(w, r)
}
