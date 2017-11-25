package user

import "net/http"

func (users *Server) List(context *Context) {
	userlist, err := context.Users.List()
	if err != nil {
		context.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	context.Data["Users"] = userlist
	context.Render("users")
}
