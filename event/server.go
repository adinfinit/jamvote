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

func (event *Server) Register(router *mux.Router) {
	router.HandleFunc("/", event.HandlerMaybe(event.List))
	router.HandleFunc("/event/create", event.HandlerMaybe(event.CreateEvent))

	router.HandleFunc("/event/{eventid}", event.Handler(event.Dashboard))
	router.HandleFunc("/event/{eventid}/edit", event.Handler(event.EditEvent))
	router.HandleFunc("/event/{eventid}/jammers", event.Handler(event.Jammers))
	router.HandleFunc("/event/{eventid}/teams", event.Handler(event.Teams))
	router.HandleFunc("/event/{eventid}/voting", event.Handler(event.Voting))
	router.HandleFunc("/event/{eventid}/fill-queue", event.Handler(event.FillQueue))
	router.HandleFunc("/event/{eventid}/progress", event.Handler(event.Progress))
	router.HandleFunc("/event/{eventid}/reveal", event.Handler(event.Reveal))
	router.HandleFunc("/event/{eventid}/results", event.Handler(event.Results))

	router.HandleFunc("/event/{eventid}/team/create", event.Handler(event.CreateTeam))
	router.HandleFunc("/event/{eventid}/team/{teamid}", event.Handler(event.Team))
	router.HandleFunc("/event/{eventid}/team/{teamid}/edit", event.Handler(event.EditTeam))
	router.HandleFunc("/event/{eventid}/vote/{teamid}", event.Handler(event.Vote))
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
	context.Render("event-dashboard")
}
