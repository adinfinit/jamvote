package user

import "net/http"

func (users *Server) List(context *Context) {
	if context.CurrentUser == nil {
		// TODO: show flash message
		context.Redirect("/", http.StatusTemporaryRedirect)
		return
	}

	userlist, err := context.Users.List()
	if err != nil {
		context.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	context.Data["Users"] = userlist
	context.Render("users")
}
