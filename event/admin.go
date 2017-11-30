package event

import (
	"fmt"
	"net/http"
	"strings"
)

func (event *Server) CreateEvent(context *Context) {
	if !context.CurrentUser.IsAdmin() {
		context.Flash("Must be admin to create events.")
		context.Redirect("/", http.StatusSeeOther)
		return
	}

	if context.Request.Method == http.MethodPost {
		if err := context.Request.ParseForm(); err != nil {
			context.FlashNow("Invalid form data: " + err.Error())
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-create")
			return
		}

		name := context.Request.FormValue("name")
		slug := context.Request.FormValue("slug")
		info := context.Request.FormValue("info")

		event := &Event{}
		event.ID = EventID(strings.ToLower(slug))
		event.Name = name
		event.Info = info

		context.Data["NewEvent"] = event

		if name == "" || !event.ID.Valid() {
			if name == "" {
				context.FlashNow("Name cannot be empty")
			}
			if !event.ID.Valid() {
				context.FlashNow("Invalid slug, can only contain 'a'-'z', '0'-'9'.")
			}

			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-create")
			return
		}

		event.Organizers = append(event.Organizers, context.CurrentUser.ID)

		err := context.Events.Create(event)
		if err != nil {
			if err == ErrExists {
				context.FlashNow(fmt.Sprintf("Event with slug %q already exists.", event.ID))
				context.Response.WriteHeader(http.StatusConflict)
			} else {
				context.FlashNow(err.Error())
				context.Response.WriteHeader(http.StatusInternalServerError)
			}
			context.Render("event-create")
			return
		}

		context.Redirect(string(event.Path()), http.StatusSeeOther)
		return
	}

	context.Render("event-create")
}

func (event *Server) EditEvent(context *Context) {
	if !context.CurrentUser.IsAdmin() {
		context.Flash("Must be admin to edit events.")
		context.Redirect("/", http.StatusSeeOther)
		return
	}

	if context.Request.Method == http.MethodPost {
		if err := context.Request.ParseForm(); err != nil {
			context.FlashNow("Invalid form data: " + err.Error())
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-edit")
			return
		}

		voting := context.Request.FormValue("voting") == "true"
		closed := context.Request.FormValue("closed") == "true"
		revealed := context.Request.FormValue("revealed") == "true"
		info := context.Request.FormValue("info")

		event := context.Event
		event.Voting = voting
		event.Closed = closed
		event.Revealed = revealed
		event.Info = info

		err := context.Events.Update(event)
		if err != nil {
			context.FlashNow(err.Error())
			context.Response.WriteHeader(http.StatusInternalServerError)
			context.Render("event-edit")
			return
		}

		context.Redirect(string(event.Path()), http.StatusSeeOther)
		return
	}

	context.Render("event-edit")
}

func (event *Server) Jammers(context *Context) {
	context.Render("todo")
	return

	if !context.CurrentUser.IsAdmin() {
		context.Flash("Must be admin to edit jammers.")
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}

	users, err := context.Users.List()
	if err != nil {
		context.FlashNow(err.Error())
	}
	context.Data["Users"] = users

	context.Render("event-jammers")
}
