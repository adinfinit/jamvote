package event

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/adinfinit/jamvote/site"
	"github.com/adinfinit/jamvote/user"
)

// CreateEvent handles page for creating a new event.
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

		var judgePercentage float64
		var err error

		if context.FormValue("judgePercentage") == "" {
			judgePercentage = 0
		} else {
			judgePercentage, err = strconv.ParseFloat(context.FormValue("judgePercentage"), 64)
		}

		if err != nil {
			context.FlashErrorNow(err.Error())
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-edit")
			return
		}

		starttime := context.FormValue("StartTime")
		endtime := context.FormValue("EndTime")

		event := &Event{}
		event.ID = EventID(strings.ToLower(slug))
		event.Name = name
		event.Theme = theme
		event.Info = info
		event.Registration = true
		event.JudgePercentage = judgePercentage

		event.Created = time.Now().UTC()

		context.Data["NewEvent"] = event

		if starttime == "" {
			event.StartTime = time.Time{}
		} else {
			t, err := time.ParseInLocation("2006-01-02T15:04", starttime, site.APTLocation)
			if err != nil {
				context.FlashErrorNow(err.Error())
				context.Response.WriteHeader(http.StatusBadRequest)
				context.Render("event-edit")
				return
			}
			event.StartTime = t
		}

		if endtime == "" {
			event.EndTime = time.Time{}
		} else {
			t, err := time.ParseInLocation("2006-01-02T15:04", endtime, site.APTLocation)
			if err != nil {
				context.FlashErrorNow(err.Error())
				context.Response.WriteHeader(http.StatusBadRequest)
				context.Render("event-edit")
				return
			}
			event.EndTime = t
		}

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

		err = context.Events.Create(event)
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

// EditEvent handles page for editing an event.
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

		var judgePercentage float64
		var err error
		if context.FormValue("judgePercentage") == "" {
			judgePercentage = 0
		} else {
			judgePercentage, err = strconv.ParseFloat(context.FormValue("judgePercentage"), 64)
		}

		if err != nil {
			context.FlashErrorNow(err.Error())
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-edit")
			return
		}

		registration := context.FormValue("registration") == "true"
		voting := context.FormValue("voting") == "true"
		closed := context.FormValue("closed") == "true"
		revealed := context.FormValue("revealed") == "true"
		info := context.FormValue("info")

		starttime := context.FormValue("StartTime")
		endtime := context.FormValue("EndTime")

		votingopens := context.FormValue("VotingOpens")
		votingcloses := context.FormValue("VotingCloses")

		event := context.Event
		event.Theme = theme
		event.JudgePercentage = judgePercentage
		event.Registration = registration
		event.Voting = voting
		event.Closed = closed
		event.Revealed = revealed
		event.Info = info

		if starttime == "" {
			event.StartTime = time.Time{}
		} else {
			t, err := time.ParseInLocation("2006-01-02T15:04", starttime, site.APTLocation)
			if err != nil {
				context.FlashErrorNow(err.Error())
				context.Response.WriteHeader(http.StatusBadRequest)
				context.Render("event-edit")
				return
			}
			event.StartTime = t
		}

		if endtime == "" {
			event.EndTime = time.Time{}
		} else {
			t, err := time.ParseInLocation("2006-01-02T15:04", endtime, site.APTLocation)
			if err != nil {
				context.FlashErrorNow(err.Error())
				context.Response.WriteHeader(http.StatusBadRequest)
				context.Render("event-edit")
				return
			}
			event.EndTime = t
		}

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
			t, err := time.ParseInLocation("2006-01-02T15:04", votingcloses, site.APTLocation)
			if err != nil {
				context.FlashErrorNow(err.Error())
				context.Response.WriteHeader(http.StatusBadRequest)
				context.Render("event-edit")
				return
			}
			event.VotingCloses = t
		}

		err = context.Events.Update(event)
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

// Jammers handles managing registered jammers for an event.
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

		jammersAdded := []user.UserID{}
		jammersRemoved := []user.UserID{}
		judgesAdded := []user.UserID{}
		judgesRemoved := []user.UserID{}

		for _, u := range users {
			jammersBefore := context.FormValue(fmt.Sprintf("%v.Jammer.Start", u.ID)) == "approved"
			jammersAfter := context.FormValue(fmt.Sprintf("%v.Jammer", u.ID)) == "approved"

			judgesBefore := context.FormValue(fmt.Sprintf("%v.Judge.Start", u.ID)) == "isjudge"
			judgesAfter := context.FormValue(fmt.Sprintf("%v.Judge", u.ID)) == "isjudge"

			if jammersBefore != jammersAfter {
				if jammersAfter {
					jammersAdded = append(jammersAdded, u.ID)
				} else {
					jammersRemoved = append(jammersRemoved, u.ID)
				}
			}

			if judgesBefore != judgesAfter {
				if judgesAfter {
					judgesAdded = append(judgesAdded, u.ID)
				} else {
					judgesRemoved = append(judgesRemoved, u.ID)
				}
			}
		}

		event := context.Event
		event.AddRemoveJammers(jammersAdded, jammersRemoved)
		event.AddRemoveJudges(judgesAdded, judgesRemoved)

		err := context.Events.Update(event)
		if err != nil {
			context.FlashErrorNow(err.Error())
			context.Response.WriteHeader(http.StatusInternalServerError)
			context.Render("event-jammers")
			return
		}

		if len(jammersRemoved) > 0 {
			context.FlashMessage(fmt.Sprintf("Removed %v jammers.", len(jammersRemoved)))
		}
		if len(judgesRemoved) > 0 {
			context.FlashMessage(fmt.Sprintf("Removed %v judges.", len(judgesRemoved)))
		}
		if len(jammersAdded) > 0 {
			context.FlashMessage(fmt.Sprintf("Added %v jammers.", len(jammersAdded)))
		}
		if len(judgesAdded) > 0 {
			context.FlashMessage(fmt.Sprintf("Added %v judges.", len(judgesAdded)))
		}

		context.Redirect(string(event.Path()), http.StatusSeeOther)
		return
	}

	context.Render("event-jammers")
}

// BallotsCSV returns all ballots for analysis.
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

	_ = writer.Write([]string{
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

		_ = writer.Write([]string{
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
