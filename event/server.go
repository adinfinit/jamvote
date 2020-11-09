package event

import (
	"fmt"
	"net/http"
	"path"
	"sort"

	"github.com/adinfinit/jamvote/site"
	"github.com/adinfinit/jamvote/user"

	"github.com/gorilla/mux"
)

// Server handles pages related to an event.
type Server struct {
	Site *site.Server
	DB   DB

	Users *user.Server
}

// Register registers all endpoints to router.
func (server *Server) Register(router *mux.Router) {
	router.HandleFunc("/", server.HandlerMaybe(server.List))
	router.HandleFunc("/event/create", server.HandlerMaybe(server.CreateEvent))

	router.HandleFunc("/event/{eventid}", server.Handler(server.Dashboard))
	router.HandleFunc("/event/{eventid}/edit", server.Handler(server.EditEvent))
	router.HandleFunc("/event/{eventid}/jammers", server.Handler(server.Jammers))
	router.HandleFunc("/event/{eventid}/linking", server.Handler(server.Linking))
	router.HandleFunc("/event/{eventid}/linking-approve-all", server.Handler(server.LinkingApproveAll))
	router.HandleFunc("/event/{eventid}/teams", server.Handler(server.Teams))
	router.HandleFunc("/event/{eventid}/voting", server.Handler(server.Voting))
	router.HandleFunc("/event/{eventid}/fill-queue", server.Handler(server.FillQueue))
	router.HandleFunc("/event/{eventid}/progress", server.Handler(server.Progress))
	router.HandleFunc("/event/{eventid}/reveal", server.Handler(server.Reveal))
	router.HandleFunc("/event/{eventid}/results", server.Handler(server.Results))

	router.HandleFunc("/event/{eventid}/ballots.csv", server.Handler(server.BallotsCSV))

	router.HandleFunc("/event/{eventid}/team/create", server.Handler(server.CreateTeam))
	router.HandleFunc("/event/{eventid}/team/{teamid}", server.Handler(server.Team))
	router.HandleFunc("/event/{eventid}/team/{teamid}/edit", server.Handler(server.EditTeam))
	router.HandleFunc("/event/{eventid}/vote/{teamid}", server.Handler(server.Vote))
}

// Path returns a proper route for an event.
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

// List lists all events.
func (server *Server) List(context *Context) {
	events, err := context.Events.List()
	if err != nil {
		context.FlashErrorNow(err.Error())
		context.Response.WriteHeader(http.StatusInternalServerError)
	}

	sort.Slice(events, func(i, k int) bool {
		return events[i].Less(events[k])
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

// Dashboard returns main page for an event.
func (server *Server) Dashboard(context *Context) {
	// TODO: deduplicate
	teams, err := context.Events.Teams(context.Event.ID)
	if err != nil {
		context.FlashErrorNow(fmt.Sprintf("Unable to get teams: %v", err))
	}

	sort.Slice(teams, func(i, k int) bool {
		return teams[i].Less(teams[k])
	})

	if context.CurrentUser != nil {
		nonsubmitted := []*Team{}
		for _, team := range teams {
			if team.HasMember(context.CurrentUser) &&
				!team.IsCompeting() {
				nonsubmitted = append(nonsubmitted, team)
			}
		}
		context.Data["NotSubmittedTeams"] = nonsubmitted
	}

	context.Render("event-dashboard")
}
