package event

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/adinfinit/jamvote/user"
)

func findUserByName(users []*user.User, name string) (*user.User, bool) {
	for _, user := range users {
		if strings.EqualFold(user.Name, name) {
			return user, true
		}
	}
	return nil, false
}
func findUserByID(users []*user.User, id user.UserID) (*user.User, bool) {
	for _, user := range users {
		if user.ID == id {
			return user, true
		}
	}
	return nil, false
}

func (server *Server) parseTeamForm(context *Context, users []*user.User) *Team {
	team := &Team{}
	team.Name = strings.TrimSpace(context.Request.FormValue("Team.Name"))

	memberNames := []string{}
	for i := 0; i < MaxTeamMembers; i++ {
		memberName := context.Request.FormValue(fmt.Sprintf("Team.Member[%v]", i))
		memberName = strings.TrimSpace(memberName)
		if memberName != "" {
			memberNames = append(memberNames, memberName)
		}
	}

	team.Members = nil
	for _, memberName := range memberNames {
		user, ok := findUserByName(users, memberName)
		if !ok {
			team.Members = append(team.Members, Member{
				ID:   0,
				Name: memberName,
			})
		} else {
			if !team.HasMember(user) {
				team.Members = append(team.Members, Member{
					ID:   user.ID,
					Name: user.Name,
				})
			}
		}
	}

	team.Game.Name = strings.TrimSpace(context.Request.FormValue("Team.Game.Name"))
	team.Game.Info = strings.TrimSpace(context.Request.FormValue("Team.Game.Info"))
	team.Game.Link.Facebook = strings.TrimSpace(context.Request.FormValue("Team.Game.Link.Facebook"))
	team.Game.Link.Jam = strings.TrimSpace(context.Request.FormValue("Team.Game.Link.Jam"))

	return team
}

func (server *Server) CreateTeam(context *Context) {
	if context.CurrentUser == nil {
		// TODO: add return address to team-creation page
		context.Flash("You must be logged in to create a team.")
		context.Redirect("/user/login", http.StatusSeeOther)
		return
	}

	if !context.CurrentUser.IsAdmin() {
		if context.Event.Closed {
			context.Flash("Event is closed and cannot be entered anymore")
			context.Redirect(context.Event.Path(), http.StatusSeeOther)
			return
		}
	}

	users, err := context.Users.List()
	if err != nil {
		context.FlashNow(fmt.Sprintf("Unable to get list of users: %v", err))
	}
	context.Data["Users"] = users
	context.Data["Team"] = &Team{}

	if context.Request.Method == http.MethodPost {
		if err := context.Request.ParseForm(); err != nil {
			context.FlashNow("Parse form: " + err.Error())
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-team-create")
			return
		}

		team := server.parseTeamForm(context, users)
		context.Data["Team"] = team
		if err := team.Verify(); err != nil {
			context.FlashNow(err.Error())
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-team-create")
			return
		}

		_, err := context.Events.CreateTeam(context.Event.ID, team)
		if err != nil {
			context.FlashNow(fmt.Sprintf("Unable to create team: %v", err))
			context.Response.WriteHeader(http.StatusInternalServerError)
			context.Render("event-team-create")
			return
		}

		context.Flash(fmt.Sprintf("Team %v created.", team.Name))
		context.Redirect(context.Event.Path("teams"), http.StatusSeeOther)
		return
	}

	context.Render("event-team-create")
}

func (server *Server) EditTeam(context *Context) {
	if !server.canEditTeam(context) {
		return
	}

	users, err := context.Users.List()
	if err != nil {
		context.FlashNow(fmt.Sprintf("Unable to get list of users: %v", err))
	}
	context.Data["Users"] = users

	if context.Request.Method == http.MethodPost {
		if err := context.Request.ParseForm(); err != nil {
			context.FlashNow("Parse form: " + err.Error())
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-team-edit")
			return
		}

		team := server.parseTeamForm(context, users)
		team.EventID = context.Team.EventID
		team.ID = context.Team.ID
		context.Data["Team"] = team

		if err := team.Verify(); err != nil {
			context.FlashNow(err.Error())
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-team-edit")
			return
		}

		err := context.Events.UpdateTeam(context.Event.ID, team)
		if err != nil {
			context.FlashNow(fmt.Sprintf("Unable to update team: %v", err))
			context.Response.WriteHeader(http.StatusInternalServerError)
			context.Render("event-team-edit")
			return
		}

		context.Redirect(context.Event.Path("teams"), http.StatusSeeOther)
		return
	}

	// Update names, if necessary
	for i, member := range context.Team.Members {
		if member.ID != 0 {
			if user, ok := findUserByID(users, member.ID); ok {
				context.Team.Members[i].Name = user.Name
			}
		}
	}

	context.Render("event-team-edit")
}

func (server *Server) Team(context *Context) {
	if context.Team == nil {
		teamid, _ := context.IntParam("teamid")
		context.Flash(fmt.Sprintf("Team %v does not exist", teamid))
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}

	context.Render("event-team")
}

func (server *Server) Teams(context *Context) {
	teams, err := context.Events.Teams(context.Event.ID)
	if err != nil {
		context.FlashNow(fmt.Sprintf("Unable to get teams: %v", err))
	}

	sort.Slice(teams, func(i, k int) bool {
		return teams[i].Name < teams[k].Name
	})

	context.Data["FullWidth"] = true
	context.Data["Teams"] = teams

	if context.CurrentUser != nil {
		yourteams := []*Team{}
		for _, team := range teams {
			if team.HasMember(context.CurrentUser) {
				yourteams = append(yourteams, team)
			}
		}
		context.Data["YourTeams"] = yourteams
	}

	context.Render("event-teams")
}

func (server *Server) canEditTeam(context *Context) bool {
	if context.Team == nil {
		teamid, _ := context.IntParam("teamid")
		context.Flash(fmt.Sprintf("Team %v does not exist.", teamid))
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return false
	}

	if !context.Team.HasEditor(context.CurrentUser) {
		context.Flash(fmt.Sprintf("You are not allowed to edit team %v.", context.Team.ID))
		context.Redirect(context.Event.Path("team", context.Team.ID.String()), http.StatusSeeOther)
		return false
	}
	return true
}
