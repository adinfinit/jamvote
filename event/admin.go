package event

import (
	"net/http"
	"strings"
	"time"
)

func (event *Server) CreateEvent(context *Context) {
	if !context.CurrentUser.IsAdmin() {
		//TODO: flash access denied
		context.Redirect("/", http.StatusTemporaryRedirect)
		return
	}

	if context.Request.Method == http.MethodPost {
		if err := context.Request.ParseForm(); err != nil {
			//TODO: flash error
			context.Error("Parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		name := context.Request.FormValue("name")
		slug := context.Request.FormValue("slug")
		info := context.Request.FormValue("info")

		event := &Event{}
		event.ID = EventID(strings.ToLower(slug))
		event.Name = name
		event.Info = info

		if name == "" || !event.ID.Valid() {
			//TODO: flash error
			context.Error("Invalid data", http.StatusBadRequest)
			return
		}

		// TODO: check valid slug

		event.Created = time.Now()
		event.Organizers = append(event.Organizers, context.CurrentUser.ID)

		err := context.Events.Create(event)
		if err != nil {
			//TODO: flash error
			context.Error(err.Error(), http.StatusBadRequest)
			return
		}

		context.Redirect(string(event.Path()), http.StatusTemporaryRedirect)
		return
	}

	context.Render("event-create")
}
