package event

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/adinfinit/jamvote/site"
	"github.com/adinfinit/jamvote/user"
)

func (server *Server) CreateEvent(context *Context) {
	if !context.CurrentUser.IsAdmin() {
		context.FlashError("Must be admin to create events.")
		context.Redirect("/", http.StatusSeeOther)
		return
	}

	if context.Request.Method == http.MethodPost {
		if err := context.Request.ParseForm(); err != nil {
			context.FlashErrorNow("Invalid form data: " + err.Error())
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-create")
			return
		}

		name := context.FormValue("name")
		theme := context.FormValue("theme")
		slug := context.FormValue("slug")
		info := context.FormValue("info")

		event := &Event{}
		event.ID = EventID(strings.ToLower(slug))
		event.Name = name
		event.Theme = theme
		event.Info = info
		event.Registration = true

		context.Data["NewEvent"] = event

		if name == "" || !event.ID.Valid() {
			if name == "" {
				context.FlashErrorNow("Name cannot be empty")
			}
			if !event.ID.Valid() {
				context.FlashErrorNow("Invalid slug, can only contain 'a'-'z', '0'-'9'.")
			}

			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-create")
			return
		}

		event.Organizers = append(event.Organizers, context.CurrentUser.ID)

		err := context.Events.Create(event)
		if err != nil {
			if err == ErrExists {
				context.FlashErrorNow(fmt.Sprintf("Event with slug %q already exists.", event.ID))
				context.Response.WriteHeader(http.StatusConflict)
			} else {
				context.FlashErrorNow(err.Error())
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

func (server *Server) EditEvent(context *Context) {
	if !context.CurrentUser.IsAdmin() {
		context.FlashError("Must be admin to edit events.")
		context.Redirect("/", http.StatusSeeOther)
		return
	}

	if context.Request.Method == http.MethodPost {
		if err := context.Request.ParseForm(); err != nil {
			context.FlashErrorNow("Invalid form data: " + err.Error())
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-edit")
			return
		}

		theme := context.FormValue("theme")
		registration := context.FormValue("registration") == "true"
		voting := context.FormValue("voting") == "true"
		closed := context.FormValue("closed") == "true"
		revealed := context.FormValue("revealed") == "true"
		info := context.FormValue("info")

		votingopens := context.FormValue("VotingOpens")
		votingcloses := context.FormValue("VotingCloses")

		event := context.Event
		event.Theme = theme
		event.Registration = registration
		event.Voting = voting
		event.Closed = closed
		event.Revealed = revealed
		event.Info = info

		if votingopens == "" {
			event.VotingOpens = time.Time{}
		} else {
			t, err := time.ParseInLocation("2006-01-02T15:04", votingopens, site.APTLocation)
			if err != nil {
				context.FlashErrorNow(err.Error())
				context.Response.WriteHeader(http.StatusBadRequest)
				context.Render("event-edit")
				return
			}
			event.VotingOpens = t
		}

		if votingcloses == "" {
			event.VotingCloses = time.Time{}
		} else {
			t, err := time.Parse("2006-01-02T15:04", votingcloses)
			if err != nil {
				context.FlashErrorNow(err.Error())
				context.Response.WriteHeader(http.StatusBadRequest)
				context.Render("event-edit")
				return
			}
			event.VotingCloses = t
		}

		err := context.Events.Update(event)
		if err != nil {
			context.FlashErrorNow(err.Error())
			context.Response.WriteHeader(http.StatusInternalServerError)
			context.Render("event-edit")
			return
		}

		context.Redirect(string(event.Path()), http.StatusSeeOther)
		return
	}

	context.Render("event-edit")
}

func (server *Server) Jammers(context *Context) {
	if !context.CurrentUser.IsAdmin() {
		context.FlashError("Must be admin to edit jammers.")
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}

	users, err := context.Users.List()
	if err != nil {
		context.FlashErrorNow(err.Error())
	}
	context.Data["Users"] = users

	if context.Request.Method == http.MethodPost {
		if err := context.Request.ParseForm(); err != nil {
			context.FlashErrorNow("Invalid form data: " + err.Error())
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-jammers")
			return
		}

		added := []user.UserID{}
		removed := []user.UserID{}

		for _, u := range users {
			before := context.FormValue(fmt.Sprintf("%v.Start", u.ID)) == "approved"
			after := context.FormValue(fmt.Sprintf("%v", u.ID)) == "approved"

			if before != after {
				if after {
					added = append(added, u.ID)
				} else {
					removed = append(removed, u.ID)
				}
			}
		}

		event := context.Event
		event.AddRemoveJammers(added, removed)

		err := context.Events.Update(event)
		if err != nil {
			context.FlashErrorNow(err.Error())
			context.Response.WriteHeader(http.StatusInternalServerError)
			context.Render("event-jammers")
			return
		}

		s := ""
		if len(removed) > 0 {
			s += fmt.Sprintf(" Removed %v jammers.", len(removed))
		}
		if len(added) > 0 {
			s += fmt.Sprintf(" Added %v jammers.", len(added))
		}
		if s != "" {
			context.FlashError(s)
		}

		context.Redirect(string(event.Path()), http.StatusSeeOther)
		return
	}

	context.Render("event-jammers")
}

func (server *Server) BallotsCSV(context *Context) {
	if !context.CurrentUser.IsAdmin() {
		context.FlashError("Must be admin to edit jammers.")
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}

	context.Response.Header().Set("Content-Type", "text/csv")

	ballots, err := context.Events.Ballots(context.Event.ID)
	if err != nil {
		context.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	teams, err := context.Events.Teams(context.Event.ID)
	if err != nil {
		context.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	teambyid := map[TeamID]*Team{}
	for _, team := range teams {
		teambyid[team.ID] = team
	}

	users, err := context.Users.List()
	if err != nil {
		context.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	userbyid := map[user.UserID]*user.User{}
	for _, user := range users {
		userbyid[user.ID] = user
	}

	writer := csv.NewWriter(context.Response)
	defer writer.Flush()

	writer.Write([]string{
		"VoterID",
		"VoterName",
		"TeamID",
		"TeamName",
		"GameName",
		"Theme",
		"Enjoyment",
		"Aesthetics",
		"Innovation",
		"Bonus",
		"Overall",
	})

	for _, ballot := range ballots {
		if !ballot.Completed {
			continue
		}
		voter := userbyid[ballot.Voter]
		team := teambyid[ballot.Team]

		writer.Write([]string{
			voter.ID.String(),
			voter.Name,
			team.ID.String(),
			team.Name,
			team.Game.Name,
			fmt.Sprintf("%.1f", ballot.Theme.Score),
			fmt.Sprintf("%.1f", ballot.Enjoyment.Score),
			fmt.Sprintf("%.1f", ballot.Aesthetics.Score),
			fmt.Sprintf("%.1f", ballot.Innovation.Score),
			fmt.Sprintf("%.1f", ballot.Bonus.Score),
			fmt.Sprintf("%.2f", ballot.Overall.Score),
		})
	}
}
