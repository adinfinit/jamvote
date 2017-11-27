package event

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (event *Server) CreateEvent(context *Context) {
	if !context.CurrentUser.IsAdmin() {
		context.Render("event-create")
		return
	}

	if context.Request.Method == http.MethodPost {
		if err := context.Request.ParseForm(); err != nil {
			context.Data["Flashes"] = []string{"Invalid form data."}
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
			flashes := []string{}
			if name == "" {
				flashes = append(flashes, "Name cannot be empty.")
			}
			if !event.ID.Valid() {
				flashes = append(flashes, "Invalid slug, can only contain 'a'-'z', '0'-'9'.")
			}
			context.FlashNow(flashes...)

			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-create")
			return
		}

		event.Created = time.Now()
		event.Stage = Started
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

		context.Redirect(string(event.Path()), http.StatusTemporaryRedirect)
		return
	}

	context.Render("event-create")
}
