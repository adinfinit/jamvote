package user

import (
	"log"
	"net/http"
	"path"

	"github.com/adinfinit/jamvote/auth"
	"github.com/gorilla/mux"
)

type Server struct {
	Auth *auth.Service
}

func (users *Server) Register(router *mux.Router) {
	router.HandleFunc("/user", users.Handler(users.RedirectToEdit))
	router.HandleFunc("/users", users.Handler(users.List))
	router.HandleFunc("/user/logged-in", users.Handler(users.LoggedIn))
	router.HandleFunc("/user/{userid}/edit", users.Handler(users.Edit))
	router.HandleFunc("/user/login", users.Handler(users.Login))
	router.HandleFunc("/user/logout", users.Handler(users.Logout))
	router.HandleFunc("/user/{userid}", users.Handler(users.Profile))
}

func getUserID(context *Context) (UserID, bool) {
	id, ok := context.IntParam("userid")
	return UserID(id), ok
}

func (users *Server) RedirectToEdit(context *Context) {
	if context.CurrentUser == nil {
		context.Redirect("/user/login", http.StatusSeeOther)
		return
	}

	userurl := path.Join("/user", context.CurrentUser.ID.String(), "edit")
	context.Redirect(userurl, http.StatusSeeOther)
}

func (users *Server) LoggedIn(context *Context) {
	if context.CurrentUser == nil || context.CurrentUser.NewUser {
		users.RedirectToEdit(context)
		return
	}

	context.Redirect("/", http.StatusSeeOther)
}

func (users *Server) Edit(context *Context) {
	userid, ok := getUserID(context)
	if !ok {
		context.Error("User ID not specified", http.StatusBadRequest)
		return
	}

	user, err := context.Users.ByID(userid)
	if err != nil {
		context.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	if !context.CurrentUser.Equals(user) && !context.CurrentUser.Admin {
		context.Flash("Editing user not allowed.")
		context.Redirect(path.Join("/user", user.ID.String()), http.StatusSeeOther)
		return
	}

	if context.Request.Method == http.MethodPost {
		if err := context.Request.ParseForm(); err != nil {
			context.Error("Parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		name := context.Request.FormValue("name")
		facebook := context.Request.FormValue("facebook")
		github := context.Request.FormValue("github")
		admin := context.Request.FormValue("admin") == "true"

		// only other admin can change admin status
		if !context.CurrentUser.Admin {
			// TODO: add flash
			admin = user.Admin
		}

		if name != user.Name ||
			facebook != user.Facebook ||
			github != user.Github ||
			admin != user.Admin {

			user.Name = name
			user.Facebook = facebook
			user.Github = github
			user.Admin = admin

			err := context.Users.Update(user)
			if err != nil {
				log.Printf("user.Edit: update %q: %v", userid, err)
				context.Error("", http.StatusInternalServerError)
				return
			}
		}
	}

	context.Data["User"] = user
	context.Render("user-edit")
}

func (users *Server) Profile(context *Context) {
	userid, ok := getUserID(context)
	if !ok {
		context.Error("User ID not specified", http.StatusBadRequest)
		return
	}

	user, err := context.Users.ByID(userid)
	if err != nil {
		context.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	context.Data["User"] = user
	context.Render("user-view")
}

func (users *Server) Login(context *Context) {
	context.Data["Logins"] = users.Auth.Links(context.Request)
	context.Render("user-login")
}

func (users *Server) Logout(context *Context) {
	users.Auth.Logout(context.Response, context.Request)
}
