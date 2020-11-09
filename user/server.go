package user

import (
	"net/http"
	"path"

	"github.com/adinfinit/jamvote/auth"
	"github.com/adinfinit/jamvote/site"

	"github.com/gorilla/mux"
)

// Server implements user management endpoints.
type Server struct {
	Site *site.Server
	DB   DB

	Auth *auth.Server
}

// Register registers endpoints to router.
func (server *Server) Register(router *mux.Router) {
	router.HandleFunc("/user", server.Handler(server.RedirectToEdit))
	router.HandleFunc("/users", server.Handler(server.List))
	router.HandleFunc("/user/logged-in", server.Handler(server.LoggedIn))
	router.HandleFunc("/user/login", server.Handler(server.Login))
	router.HandleFunc("/user/logout", server.Handler(server.Logout))
}

// RedirectToEdit redirects request to edit page.
func (server *Server) RedirectToEdit(context *Context) {
	if context.CurrentUser == nil {
		context.Redirect("/user/login", http.StatusSeeOther)
		return
	}

	userurl := path.Join("/user", context.CurrentUser.ID.String(), "edit")
	context.Redirect(userurl, http.StatusSeeOther)
}

// LoggedIn is called after user has logged in.
func (server *Server) LoggedIn(context *Context) {
	if context.CurrentUser == nil || context.CurrentUser.NewUser {
		context.FlashMessage("Please update your full name.")
		server.RedirectToEdit(context)
		return
	}

	context.Redirect("/", http.StatusSeeOther)
}

// Login renders login page.
func (server *Server) Login(context *Context) {
	context.Data["Logins"] = server.Auth.Links(context.Request)
	context.Render("user-login")
}

// Logout logs out.
func (server *Server) Logout(context *Context) {
	server.Auth.Logout(context.Response, context.Request)
}
