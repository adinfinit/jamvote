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

type Server struct {
	Site  *site.Server
	Users *user.Server
}

func (server *Server) Register(router *mux.Router) {
	router.HandleFunc("/", server.HandlerMaybe(server.List))
	router.HandleFunc("/event/create", server.HandlerMaybe(server.CreateEvent))

	router.HandleFunc("/event/{eventid}", server.Handler(server.Dashboard))
	router.HandleFunc("/event/{eventid}/edit", server.Handler(server.EditEvent))
	router.HandleFunc("/event/{eventid}/jammers", server.Handler(server.Jammers))
	router.HandleFunc("/event/{eventid}/teams", server.Handler(server.Teams))
	router.HandleFunc("/event/{eventid}/voting", server.Handler(server.Voting))
	router.HandleFunc("/event/{eventid}/fill-queue", server.Handler(server.FillQueue))
	router.HandleFunc("/event/{eventid}/progress", server.Handler(server.Progress))
	router.HandleFunc("/event/{eventid}/reveal", server.Handler(server.Reveal))
	router.HandleFunc("/event/{eventid}/results", server.Handler(server.Results))

	router.HandleFunc("/event/{eventid}/team/create", server.Handler(server.CreateTeam))
	router.HandleFunc("/event/{eventid}/team/{teamid}", server.Handler(server.Team))
	router.HandleFunc("/event/{eventid}/team/{teamid}/edit", server.Handler(server.EditTeam))
	router.HandleFunc("/event/{eventid}/vote/{teamid}", server.Handler(server.Vote))
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

func (server *Server) List(context *Context) {
	events, err := context.Events.List()
	if err != nil {
		context.FlashErrorNow(err.Error())
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

func (server *Server) Dashboard(context *Context) {
	context.Render("event-dashboard")
}
