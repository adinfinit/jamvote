package user

import (
	"log"
	"net/http"

	"github.com/adinfinit/rater/site"
	"github.com/gorilla/mux"

	"google.golang.org/appengine"
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

func (users *Server) EditProfile(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	if u == nil {
		http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)
		return
	}

	user := &User{
		ID:    ID(u.ID),
		Name:  u.Email,
		Email: u.Email,
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
