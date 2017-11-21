package user

import (
	"log"
	"net/http"

	"github.com/adinfinit/jamvote/site"
	"github.com/gorilla/mux"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

type Server struct {
	Renderer *site.Renderer
}

func (users *Server) Register(router *mux.Router) {
	router.HandleFunc("/user", users.EditProfile) // REDIRECT TO SELF
	router.HandleFunc("/user/login", users.Login)
	router.HandleFunc("/user/logout", users.Logout)
	router.HandleFunc("/user/{userid}", users.Profile)
}

type GoogleAuth struct {
	UserID ID
}

func (users *Server) EditProfile(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	googleuser := user.Current(c)

	if googleuser == nil {
		http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)
		return
	}

	authkey := datastore.NewKey(c, "Auth", "google-"+string(googleuser.ID), 0, nil)

	auth := &GoogleAuth{}
	err := datastore.Get(c, authkey, auth)

	if err == datastore.ErrNoSuchEntity {
		// create new user
		siteuser := &User{
			ID:    ID(googleuser.ID),
			Name:  googleuser.Email,
			Email: googleuser.Email,
		}

		// create new user
		userkey := datastore.NewKey(c, "User", string(siteuser.ID), 0, nil)
		if _, err := datastore.Put(c, userkey, siteuser); err != nil {
			http.Error(w, "Create new user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		auth.UserID = siteuser.ID

		// add authentication mapping
		if _, err := datastore.Put(c, authkey, auth); err != nil {
			http.Error(w, "Create new auth: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	siteuser := &User{}
	userkey := datastore.NewKey(c, "User", string(auth.UserID), 0, nil)

	if err := datastore.Get(c, userkey, siteuser); err != nil {
		http.Error(w, "Get user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	siteuser.ID = auth.UserID

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		facebook := r.FormValue("facebook")

		if name != siteuser.Name || email != siteuser.Email || facebook != siteuser.Facebook {
			siteuser.Name = name
			siteuser.Email = email
			siteuser.Facebook = facebook

			if _, err := datastore.Put(c, userkey, siteuser); err != nil {
				http.Error(w, "Put user: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	users.Renderer.Render(w, "user-edit", map[string]interface{}{
		"User": siteuser,
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
	c := appengine.NewContext(r)

	googlelogin, err := user.LoginURL(c, "/user")
	if err != nil {
		log.Println(err)
	}

	users.Renderer.Render(w, "user-login", map[string]interface{}{
		"Logins": []Login{{"Google", googlelogin}},
	})
}

func (users *Server) Logout(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	googlelogout, err := user.LogoutURL(c, "/")
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, googlelogout, http.StatusTemporaryRedirect)
}
