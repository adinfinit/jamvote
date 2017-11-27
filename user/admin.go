package user

import "net/http"

func (users *Server) List(context *Context) {
	if context.CurrentUser == nil {
		context.Flash("Must be logged in to see users.")
		context.Redirect("/", http.StatusSeeOther)
		return
	}

	userlist, err := context.Users.List()
	if err != nil {
		context.FlashNow(err.Error())
	}

	context.Data["Users"] = userlist
	context.Render("users")
}
