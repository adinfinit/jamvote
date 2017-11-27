package event

import "net/http"

func (event *Server) Voting(context *Context) {
	if !context.CurrentUser.IsAdmin() {
		if !context.Event.Voting {
			context.Flash("Voting has not yet started.")
			context.Redirect(context.Event.Path(), http.StatusSeeOther)
			return
		}
	}

	if context.Event.Closed {
		context.FlashNow("Voting has been closed.")
	}

	context.Render("todo")
}

func (event *Server) Results(context *Context) {
	if !context.CurrentUser.IsAdmin() {
		if !context.Event.Revealed {
			context.Flash("Results have not yet been revealed.")
			context.Redirect(context.Event.Path(), http.StatusSeeOther)
			return
		}
	}

	context.Render("todo")
}
