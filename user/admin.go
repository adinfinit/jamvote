package user

import "net/http"

func (server *Server) List(context *Context) {
	if context.CurrentUser == nil {
		context.FlashError("Must be logged in to see users.")
		context.Redirect("/", http.StatusSeeOther)
		return
	}

	userlist, err := context.Users.List()
	if err != nil {
		context.FlashErrorNow(err.Error())
	}

	context.Data["Users"] = userlist
	context.Render("users")
}
