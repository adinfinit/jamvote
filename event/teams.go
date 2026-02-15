package event

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/adinfinit/jamvote/user"
)

// findUserByName finds a user by name from users.
func findUserByName(users []*user.User, name string) (*user.User, bool) {
	for _, user := range users {
		if strings.EqualFold(user.Name, name) {
			return user, true
		}
	}
	return nil, false
}

// findUserByID finds a user by ID from users.
func findUserByID(users []*user.User, id user.UserID) (*user.User, bool) {
	for _, user := range users {
		if user.ID == id {
			return user, true
		}
	}
	return nil, false
}

// parseTeamForm parses edited team page.
func (server *Server) parseTeamForm(context *Context, users []*user.User) *Team {
	team := &Team{}
	team.Name = context.FormValue("Team.Name")

	memberNames := []string{}
	for i := range MaxTeamMembers {
		memberName := context.FormValue(fmt.Sprintf("Team.Member[%v]", i))
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

	team.Game.Name = context.FormValue("Team.Game.Name")
	team.Game.Info = context.FormValue("Team.Game.Info")
	team.Game.Noncompeting = context.FormValue("Team.Game.Noncompeting") == "true"
	team.Game.Link.Jam = context.FormValue("Team.Game.Link.Jam")
	team.Game.Link.Download = context.FormValue("Team.Game.Link.Download")
	team.Game.Link.Facebook = context.FormValue("Team.Game.Link.Facebook")

	return team
}

// CreateTeam handles page for creating a new team.
func (server *Server) CreateTeam(context *Context) {
	if context.CurrentUser == nil {
		// TODO: add return address to team-creation page
		context.FlashError("You must be logged in to create a team.")
		context.Redirect("/user/login", http.StatusSeeOther)
		return
	}

	if !context.Event.CanRegister(context.CurrentUser) {
		context.FlashMessage("Registration is closed and cannot be entered anymore")
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}

	users, err := context.Users.List()
	if err != nil {
		context.FlashErrorNow(fmt.Sprintf("Unable to get list of users: %v", err))
	}
	context.Data["Users"] = users
	context.Data["Team"] = &Team{}

	if context.Request.Method == http.MethodPost {
		if err := context.Request.ParseForm(); err != nil {
			context.FlashErrorNow("Parse form: " + err.Error())
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-team-create")
			return
		}

		team := server.parseTeamForm(context, users)
		context.Data["Team"] = team
		if err := team.Verify(); err != nil {
			context.FlashErrorNow(err.Error())
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-team-create")
			return
		}

		_, err := context.Events.CreateTeam(context.Event.ID, team)
		if err != nil {
			context.FlashErrorNow(fmt.Sprintf("Unable to create team: %v", err))
			context.Response.WriteHeader(http.StatusInternalServerError)
			context.Render("event-team-create")
			return
		}

		context.FlashError(fmt.Sprintf("Team %v created.", team.Name))
		context.Redirect(context.Event.Path("teams"), http.StatusSeeOther)
		return
	}

	context.Render("event-team-create")
}

// EditTeam handles page for editing a team.
func (server *Server) EditTeam(context *Context) {
	if !server.canEditTeam(context) {
		return
	}

	users, err := context.Users.List()
	if err != nil {
		context.FlashErrorNow(fmt.Sprintf("Unable to get list of users: %v", err))
	}
	context.Data["Users"] = users

	if context.Request.Method == http.MethodPost {
		if err := context.Request.ParseForm(); err != nil {
			context.FlashErrorNow("Parse form: " + err.Error())
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-team-edit")
			return
		}

		team := server.parseTeamForm(context, users)
		team.EventID = context.Team.EventID
		team.ID = context.Team.ID
		context.Data["Team"] = team

		if err := team.Verify(); err != nil {
			context.FlashErrorNow(err.Error())
			context.Response.WriteHeader(http.StatusBadRequest)
			context.Render("event-team-edit")
			return
		}

		err := context.Events.UpdateTeam(context.Event.ID, team)
		if err != nil {
			context.FlashErrorNow(fmt.Sprintf("Unable to update team: %v", err))
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

// DeleteTeam handles page for deleting a team.
func (server *Server) DeleteTeam(context *Context) {
	if !server.canDeleteTeam(context) {
		return
	}

	err := context.Events.DeleteTeam(context.Event.ID, context.Team.ID)
	if err != nil {
		context.FlashError(fmt.Sprintf("Unable to update team: %v", err))
		context.Response.WriteHeader(http.StatusSeeOther)
		context.Redirect(context.Event.Path("team", context.Team.ID.String()), http.StatusSeeOther)
		return
	}

	context.FlashMessage(fmt.Sprintf("Team %v deleted.", context.Team.ID))
	context.Redirect(context.Event.Path("teams"), http.StatusSeeOther)
}

// Team displays team information.
func (server *Server) Team(context *Context) {
	if context.Team == nil {
		teamid, _ := context.IntParam("teamid")
		context.FlashError(fmt.Sprintf("Team %v does not exist", teamid))
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return
	}

	if context.Event.Revealed {
		ballots, err := context.Events.TeamBallots(context.Event.ID, context.Team.ID)
		if err != nil {
			context.FlashError(err.Error())
		}

		var aspectsInfo AspectsInfo
		for _, ballot := range ballots {
			if !ballot.Completed {
				continue
			}
			if context.CurrentUser != nil && ballot.Voter == context.CurrentUser.ID {
				context.Data["CurrentUserBallot"] = ballot
			}
			aspectsInfo.Add(&ballot.Aspects, context.Team.HasMemberID(ballot.Voter))
		}
		context.Data["Aspects"] = AspectDescriptionsWithOverall
		context.Data["AspectsInfo"] = &aspectsInfo
	}

	context.Render("event-team")
}

// Teams displays all teams.
func (server *Server) Teams(context *Context) {
	teams, err := context.Events.Teams(context.Event.ID)
	if err != nil {
		context.FlashErrorNow(fmt.Sprintf("Unable to get teams: %v", err))
	}

	sort.Slice(teams, func(i, k int) bool {
		return teams[i].Less(teams[k])
	})

	context.Data["FullWidth"] = true
	context.Data["Teams"] = teams

	maxMembers := SuggestedTeamMembers
	for _, t := range teams {
		if len(t.Members) > maxMembers {
			maxMembers = len(t.Members)
		}
	}
	context.Data["MaxMemberCount"] = maxMembers

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

// Linking displays information about users who have not been linked to their users.
func (server *Server) Linking(context *Context) {
	if !context.CurrentUser.IsAdmin() {
		context.FlashError("Must be admin to view linking information.")
		context.Redirect("/", http.StatusSeeOther)
		return
	}

	users, err := context.Users.List()
	if err != nil {
		context.FlashErrorNow(fmt.Sprintf("Unable to get users: %v", err))
	}

	teams, err := context.Events.Teams(context.Event.ID)
	if err != nil {
		context.FlashErrorNow(fmt.Sprintf("Unable to get teams: %v", err))
	}

	sort.Slice(teams, func(i, k int) bool {
		return teams[i].Less(teams[k])
	})

	type UserLink struct {
		Member Member
		User   *user.User
	}

	type Linking struct {
		Team         *Team
		Unlinked     []UserLink
		Unregistered []Member
		Unapproved   []*user.User
	}

	var linking []Linking

	for _, team := range teams {
		link := Linking{}
		link.Team = team
		for _, member := range team.Members {
			if member.ID != 0 {
				user, ok := findUserByID(users, member.ID)
				if ok && !context.Event.HasJammer(user) {
					link.Unapproved = append(link.Unapproved, user)
				}
				continue
			}

			if match, ok := findUserByName(users, member.Name); ok {
				link.Unlinked = append(link.Unlinked, UserLink{
					Member: member,
					User:   match,
				})
			} else {
				link.Unregistered = append(link.Unregistered, member)
			}
		}

		if len(link.Unlinked) > 0 || len(link.Unregistered) > 0 || len(link.Unapproved) > 0 {
			linking = append(linking, link)
		}
	}

	context.Data["Linking"] = linking

	context.Render("event-linking")
}

// LinkingApproveAll tries to associate all team members with appropriate teams and users.
func (server *Server) LinkingApproveAll(context *Context) {
	if !context.CurrentUser.IsAdmin() {
		context.FlashError("Must be admin to approve all.")
		context.Redirect("/", http.StatusSeeOther)
		return
	}

	users, err := context.Users.List()
	if err != nil {
		context.FlashError(fmt.Sprintf("Unable to get users: %v", err))
		context.Redirect(context.Event.Path("linking"), http.StatusSeeOther)
		return
	}

	teams, err := context.Events.Teams(context.Event.ID)
	if err != nil {
		context.FlashError(fmt.Sprintf("Unable to get teams: %v", err))
		context.Redirect(context.Event.Path("linking"), http.StatusSeeOther)
		return
	}

	var unapproved []user.UserID
	for _, team := range teams {
		for _, member := range team.Members {
			if member.ID != 0 {
				user, ok := findUserByID(users, member.ID)
				if ok && !context.Event.HasJammer(user) {
					unapproved = append(unapproved, user.ID)
				}
				continue
			}
		}
	}

	event := context.Event
	event.AddRemoveJammers(unapproved, nil)

	err = context.Events.Update(event)
	if err != nil {
		context.FlashError(err.Error())
	} else {
		context.FlashMessage(fmt.Sprintf("Added %v jammers.", len(unapproved)))
	}

	context.Redirect(context.Event.Path("linking"), http.StatusSeeOther)
}

// canEditTeam checks whether caller can edit the team.
func (server *Server) canEditTeam(context *Context) bool {
	if context.Team == nil {
		teamid, _ := context.IntParam("teamid")
		context.FlashError(fmt.Sprintf("Team %v does not exist.", teamid))
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return false
	}

	if !context.Team.HasEditor(context.CurrentUser) {
		context.FlashError(fmt.Sprintf("You are not allowed to edit team %v.", context.Team.ID))
		context.Redirect(context.Event.Path("team", context.Team.ID.String()), http.StatusSeeOther)
		return false
	}
	return true
}

// canDeleteTeam checks whether caller can edit the team.
func (server *Server) canDeleteTeam(context *Context) bool {
	if context.Team == nil {
		teamid, _ := context.IntParam("teamid")
		context.FlashError(fmt.Sprintf("Team %v does not exist.", teamid))
		context.Redirect(context.Event.Path(), http.StatusSeeOther)
		return false
	}

	if context.Event.Voting || context.Event.Closed || context.Event.Revealed {
		context.FlashError("Team can only be deleted during registration.")
		context.Redirect(context.Event.Path("team", context.Team.ID.String()), http.StatusSeeOther)
		return false
	}

	if !context.CurrentUser.IsAdmin() {
		context.FlashError(fmt.Sprintf("You are not allowed to delete team %v.", context.Team.ID))
		context.Redirect(context.Event.Path("team", context.Team.ID.String()), http.StatusSeeOther)
		return false
	}
	return true
}
