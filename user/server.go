package user

import (
	"net/http"

	"github.com/adinfinit/rater/site"
	"github.com/gorilla/mux"
)

type Server struct {
	Renderer *site.Renderer
}

func (users *Server) Register(router *mux.Router) {
	router.HandleFunc("/user", users.RedirectSelf) // REDIRECT TO SELF
	router.HandleFunc("/user/login", users.Login)
	router.HandleFunc("/user/logout", users.Logout)
	router.HandleFunc("/user/{userid}", users.Profile)
}

func (users *Server) RedirectSelf(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/user/self", http.StatusTemporaryRedirect)
}

func (users *Server) Profile(w http.ResponseWriter, r *http.Request) {
	userid := mux.Vars(r)["userid"]

	users.Renderer.Render(w, "user-profile", map[string]interface{}{
		"UserID": userid,
	})
}

func (users *Server) Login(w http.ResponseWriter, r *http.Request) {
	users.Renderer.Render(w, "user-login", map[string]interface{}{})
}

func (users *Server) Logout(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
