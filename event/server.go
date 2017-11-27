package event

import (
	"fmt"
	"net/http"
	"path"
	"sort"

	"github.com/adinfinit/jamvote/user"
	"github.com/gorilla/mux"
)

type Server struct {
	Users *user.Server
}

func NewServer(users *user.Server) *Server {
	return &Server{users}
}

func (event *Server) Register(router *mux.Router) {
	router.HandleFunc("/", event.HandlerMaybe(event.List))
	router.HandleFunc("/event/create", event.HandlerMaybe(event.CreateEvent))

	router.HandleFunc("/event/{eventid}", event.Handler(event.Dashboard))
	router.HandleFunc("/event/{eventid}/edit", event.Handler(event.EditEvent))
	router.HandleFunc("/event/{eventid}/voting", event.Handler(event.Voting))
	router.HandleFunc("/event/{eventid}/results", event.Handler(event.Results))

	router.HandleFunc("/event/{eventid}/team/create", event.Handler(event.CreateTeam))
	router.HandleFunc("/event/{eventid}/team/{teamid}", event.Handler(event.Team))
	router.HandleFunc("/event/{eventid}/team/{teamid}/edit", event.Handler(event.EditTeam))
	router.HandleFunc("/event/{eventid}/team/{teamid}/edit-game", event.Handler(event.EditTeamGame))
	router.HandleFunc("/event/{eventid}/vote/{teamid}", event.Handler(event.Dashboard))
}

func (event *Event) Path(subroutes ...interface{}) string {
	route := []string{"/event", string(event.ID)}
	for _, r := range subroutes {
		switch x := r.(type) {
		case int, int64, int32, uint, uint64, uint32, string:
			route = append(route, fmt.Sprint(x))
		case fmt.Stringer:
			route = append(route, x.String())
		default:
			panic("unsupported type")
		}
	}
	return path.Join(route...)
}

func (dashboard *Server) List(context *Context) {
	events, err := context.Events.List()
	if err != nil {
		context.FlashNow(err.Error())
		context.Response.WriteHeader(http.StatusInternalServerError)
	}

	sort.Slice(events, func(i, k int) bool {
		return events[i].Name < events[k].Name
	})

	byStage := struct {
		All      []*Event
		Started  []*Event
		Finished []*Event
	}{}
	byStage.All = events

	for _, event := range events {
		if !event.Closed {
			byStage.Started = append(byStage.Started, event)
		} else {
			byStage.Finished = append(byStage.Finished, event)
		}
	}

	context.Data["Events"] = byStage
	context.Render("event-list")
}

func (event *Server) Dashboard(context *Context) {
	teams, err := context.Events.Teams(context.Event.ID)
	if err != nil {
		context.FlashNow(fmt.Sprintf("Unable to get teams: %v", err))
	}

	sort.Slice(teams, func(i, k int) bool {
		return teams[i].Name < teams[k].Name
	})
	context.Data["Teams"] = teams

	if context.CurrentUser != nil {
		yourteams := []*Team{}
		otherteams := []*Team{}
		for _, team := range teams {
			if team.HasMember(context.CurrentUser) {
				yourteams = append(yourteams, team)
			} else {
				otherteams = append(otherteams, team)
			}
		}
		context.Data["YourTeams"] = yourteams
		context.Data["Teams"] = otherteams
	}

	context.Render("event-dashboard")
}
