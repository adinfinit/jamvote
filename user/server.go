package user

import (
	"context"
	"log"
	"net/http"
	"path"
	"strconv"

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
	router.HandleFunc("/user", users.Redirect)
	router.HandleFunc("/user/{userid}/edit", users.Edit) // REDIRECT TO SELF
	router.HandleFunc("/user/login", users.Login)
	router.HandleFunc("/user/logout", users.Logout)
	router.HandleFunc("/user/{userid}", users.Profile)
}

type CredentialsUser struct {
	UserKey *datastore.Key

	Provider string `datastore:",noindex"`
	Email    string `datastore:",noindex"`
	Name     string `datastore:",noindex"`
}

func getUserID(r *http.Request) (int64, bool) {
	s := mux.Vars(r)["userid"]
	if s == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}

func (users *Server) Current(c context.Context, r *http.Request) *User {
	cred := users.Auth.CurrentCredentials(r)
	if cred == nil {
		return nil
	}

	authkey := datastore.NewKey(c, "Auth", cred.ID, 0, nil)
	creduser := &CredentialsUser{}
	err := datastore.Get(c, authkey, creduser)

	if err == datastore.ErrNoSuchEntity {
		// create new user
		user := &User{
			Name: cred.Name,
		}

		incompletekey := datastore.NewIncompleteKey(c, "User", nil)
		userkey, err := datastore.Put(c, incompletekey, user)
		if err != nil {
			log.Printf("Create user: %v", err)
			return nil
		}

		creduser.UserKey = userkey
		creduser.Provider = cred.Provider
		creduser.Email = cred.Email
		creduser.Name = cred.Name

		_, err = datastore.Put(c, authkey, creduser)
		if err != nil {
			log.Printf("Create user credentials: %v", err)
			return nil
		}
	} else if err != nil {
		log.Printf("Get auth failed: %v", err)
		return nil
	}

	user := &User{}
	if err := datastore.Get(c, creduser.UserKey, user); err != nil {
		log.Printf("Load user: %v", err)
		return nil
	}
	user.ID = ID(creduser.UserKey.IntID())

	return user
}

func (users *Server) UserByID(c context.Context, id int64) *User {
	user := &User{}
	key := datastore.NewKey(c, "User", "", id, nil)
	if err := datastore.Get(c, key, user); err != nil {
		return nil
	}
	user.ID = ID(key.IntID())
	return user
}

func (users *Server) UserByKey(c context.Context, key *datastore.Key) *User {
	user := &User{}
	if err := datastore.Get(c, key, user); err != nil {
		return nil
	}
	user.ID = ID(key.IntID())
	return user
}

func (users *Server) Redirect(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	currentuser := users.Current(c, r)
	if currentuser == nil {
		http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)
		return
	}

	userurl := path.Join("/user", currentuser.ID.String(), "edit")
	http.Redirect(w, r, userurl, http.StatusTemporaryRedirect)
}

func (users *Server) Edit(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	userid, ok := getUserID(r)
	if !ok {
		http.Error(w, "User ID not specified", http.StatusBadRequest)
		return
	}

	currentuser := users.Current(c, r)
	user := users.UserByID(c, int64(userid))
	if user == nil {
		http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		facebook := r.FormValue("facebook")
		github := r.FormValue("github")

		if name != user.Name ||
			facebook != user.Facebook ||
			github != user.Github {

			user.Name = name
			user.Facebook = facebook
			user.Github = github

			c := appengine.NewContext(r)
			userkey := datastore.NewKey(c, "User", "", int64(user.ID), nil)
			if _, err := datastore.Put(c, userkey, user); err != nil {
				http.Error(w, "Put user: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	users.Renderer.Render(w, "user-edit", map[string]interface{}{
		"CurrentUser": currentuser,
		"User":        user,
	})
}

func (users *Server) Profile(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	currentuser := users.Current(c, r)

	userid, ok := getUserID(r)
	if !ok {
		http.Error(w, "User ID not specified", http.StatusBadRequest)
		return
	}
	user := users.UserByID(c, userid)

	users.Renderer.Render(w, "user-show", map[string]interface{}{
		"CurrentUser": currentuser,
		"User":        user,
	})
}

func (users *Server) Login(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	currentuser := users.Current(c, r)

	users.Renderer.Render(w, "user-login", map[string]interface{}{
		"CurrentUser": currentuser,
		"Logins":      users.Auth.Links(r),
	})
}

func (users *Server) Logout(w http.ResponseWriter, r *http.Request) {
	users.Auth.Logout(w, r)
}
