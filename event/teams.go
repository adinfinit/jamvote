package event

import (
	"fmt"
	"net/http"
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

func (server *Server) readTeamInfo(context *Context, users []*user.User) (*Team, bool) {
	if err := context.Request.ParseForm(); err != nil {
		context.FlashNow("Parse form: " + err.Error())
		return nil, false
	}

	teamname := strings.TrimSpace(context.Request.FormValue("name"))
	if teamname == "" {
		context.FlashNow("Team name cannot be empty.")
		return nil, false
	}

	memberNames := []string{}
	for i := 0; i < 5; i++ {
		memberName := context.Request.FormValue(fmt.Sprintf("member[%v]", i))
		memberName = strings.TrimSpace(memberName)
		if memberName != "" {
			memberNames = append(memberNames, memberName)
		}
	}

	team := &Team{}
	team.Name = teamname
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

	return team, true
}

func (event *Server) CreateTeam(context *Context) {
	if context.CurrentUser == nil {
		// TODO: add return address to team-creation page
		context.Flash("You must be logged in to create a team.")
		context.Redirect("/user/login", http.StatusTemporaryRedirect)
		return
	}

	users, err := context.Users.List()
	if err != nil {
		context.FlashNow(fmt.Sprintf("Unable to get list of users: %v", err))
	}
	context.Data["Users"] = users

	if context.Request.Method == http.MethodPost {
		team, ok := event.readTeamInfo(context, users)
		if !ok {
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-team-create")
			return
		}

		if !team.HasMember(context.CurrentUser) {
			team.Members = append([]Member{{
				ID:   context.CurrentUser.ID,
				Name: context.CurrentUser.Name,
			}}, team.Members...)
		}

		teamid, err := context.Events.CreateTeam(context.Event.ID, team)
		if err != nil {
			context.FlashNow(fmt.Sprintf("Unable to create team: %v", err))
			context.Response.WriteHeader(http.StatusInternalServerError)
			context.Render("event-team-create")
			return
		}

		context.Redirect(context.Event.Path("team", teamid.String()), http.StatusTemporaryRedirect)
		return
	}

	context.Render("event-team-create")
}

func (event *Server) Team(context *Context) {
	if context.Team == nil {
		teamid, _ := context.IntParam("teamid")
		context.Flash(fmt.Sprintf("Team %v does not exist", teamid))
		context.Redirect(context.Event.Path(), http.StatusTemporaryRedirect)
		return
	}

	context.Render("event-team")
}

func (event *Server) EditTeam(context *Context) {
	if context.Team == nil {
		teamid, _ := context.IntParam("teamid")
		context.Flash(fmt.Sprintf("Team %v does not exist.", teamid))
		context.Redirect(context.Event.Path(), http.StatusTemporaryRedirect)
		return
	}

	if !context.Team.HasEditor(context.CurrentUser) {
		context.Flash(fmt.Sprintf("You are not allowed to edit team %v.", context.Team.ID))
		context.Redirect(context.Event.Path("team", context.Team.ID.String()), http.StatusTemporaryRedirect)
		return
	}

	users, err := context.Users.List()
	if err != nil {
		context.FlashNow(fmt.Sprintf("Unable to get list of users: %v", err))
	}

	if context.Request.Method == http.MethodPost {
		team, ok := event.readTeamInfo(context, users)
		if !ok {
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-team-edit")
			return
		}

		team.EventID = context.Team.EventID
		team.ID = context.Team.ID
		team.Entry = context.Team.Entry

		err := context.Events.UpdateTeam(context.Event.ID, team)
		if err != nil {
			context.FlashNow(fmt.Sprintf("Unable to update team: %v", err))
			context.Response.WriteHeader(http.StatusInternalServerError)
			context.Render("event-team-edit")
			return
		}

		context.Redirect(context.Event.Path("team", team.ID.String()), http.StatusTemporaryRedirect)
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

	context.Data["Users"] = users

	context.Render("event-team-edit")
}
