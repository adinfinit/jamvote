package profile

import (
	"net/http"
	"path"

	"github.com/adinfinit/jamvote/site"
	"github.com/adinfinit/jamvote/user"

	"github.com/gorilla/mux"
)

type Server struct {
	Site  *site.Server
	Users *user.Server
}

func (server *Server) Register(router *mux.Router) {
	router.HandleFunc("/user/{userid}/edit", server.Handler(server.Edit))
	router.HandleFunc("/user/{userid}", server.Handler(server.Profile))
}

func getUserID(context *Context) (user.UserID, bool) {
	id, ok := context.IntParam("userid")
	return user.UserID(id), ok
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

	teams, err := context.Teams.TeamsByUser(userid)
	if err != nil {
		context.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	context.Data["User"] = user
	context.Data["Teams"] = teams
	context.Render("user-view")
}
