package user

import (
	"net/http"
	"path"

	"github.com/adinfinit/jamvote/auth"
	"github.com/adinfinit/jamvote/site"

	"github.com/gorilla/mux"
)

type Server struct {
	Site *site.Server
	Auth *auth.Service
}

func (server *Server) Register(router *mux.Router) {
	router.HandleFunc("/user", server.Handler(server.RedirectToEdit))
	router.HandleFunc("/users", server.Handler(server.List))
	router.HandleFunc("/user/logged-in", server.Handler(server.LoggedIn))
	router.HandleFunc("/user/{userid}/edit", server.Handler(server.Edit))
	router.HandleFunc("/user/login", server.Handler(server.Login))
	router.HandleFunc("/user/logout", server.Handler(server.Logout))
	router.HandleFunc("/user/{userid}", server.Handler(server.Profile))
}

func getUserID(context *Context) (UserID, bool) {
	id, ok := context.IntParam("userid")
	return UserID(id), ok
}

func (server *Server) RedirectToEdit(context *Context) {
	if context.CurrentUser == nil {
		context.Redirect("/user/login", http.StatusSeeOther)
		return
	}

	userurl := path.Join("/user", context.CurrentUser.ID.String(), "edit")
	context.Redirect(userurl, http.StatusSeeOther)
}

func (server *Server) LoggedIn(context *Context) {
	if context.CurrentUser == nil || context.CurrentUser.NewUser {
		context.FlashMessage("Please update your full name.")
		server.RedirectToEdit(context)
		return
	}

	context.Redirect("/", http.StatusSeeOther)
}

func (server *Server) Edit(context *Context) {
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

	if !user.HasEditor(context.CurrentUser) {
		context.FlashError("Editing user not allowed.")
		context.Redirect(path.Join("/user", user.ID.String()), http.StatusSeeOther)
		return
	}

	if context.Request.Method == http.MethodPost {
		if err := context.Request.ParseForm(); err != nil {
			context.Error("Parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		user.Name = context.FormValue("name")
		user.Email = context.FormValue("email")
		user.Facebook = context.FormValue("facebook")
		user.Github = context.FormValue("github")

		admin := context.FormValue("admin") == "true"
		// only other admin can change admin status
		if context.CurrentUser.IsAdmin() {
			user.Admin = admin
		}

		err := context.Users.Update(user)
		if err != nil {
			context.FlashError(err.Error())
		} else {
			context.FlashMessage("User updated.")
		}

		context.Redirect(path.Join("/user", user.ID.String()), http.StatusSeeOther)
		return
	}

	context.Data["User"] = user
	context.Render("user-edit")
}

func (server *Server) Profile(context *Context) {
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

func (server *Server) Login(context *Context) {
	context.Data["Logins"] = server.Auth.Links(context.Request)
	context.Render("user-login")
}

func (server *Server) Logout(context *Context) {
	server.Auth.Logout(context.Response, context.Request)
}
