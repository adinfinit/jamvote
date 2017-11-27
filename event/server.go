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
	router.HandleFunc("/event/{eventid}/voting", event.Handler(event.Dashboard))
	router.HandleFunc("/event/{eventid}/results", event.Handler(event.Dashboard))

	router.HandleFunc("/event/{eventid}/team/create", event.Handler(event.CreateTeam))
	router.HandleFunc("/event/{eventid}/team/{teamid}", event.Handler(event.Team))
	router.HandleFunc("/event/{eventid}/team/{teamid}/edit", event.Handler(event.EditTeam))
	router.HandleFunc("/event/{eventid}/team/{teamid}/edit-entry", event.Handler(event.EditTeamEntry))
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
		return events[i].Created.After(events[k].Created)
	})

	byStage := struct {
		All      []*Event
		Started  []*Event
		Voting   []*Event
		Finished []*Event
	}{}
	byStage.All = events

	for _, event := range events {
		switch event.Stage {
		case Draft:
			// byStage.Draft = append(byStage.Draft, event)
		case Started:
			byStage.Started = append(byStage.Started, event)
		case Voting:
			byStage.Voting = append(byStage.Voting, event)
		case Finished:
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

	if context.CurrentUser != nil {
		yourteams := []*Team{}
		for _, team := range teams {
			if team.HasMember(context.CurrentUser) {
				yourteams = append(yourteams, team)
			}
		}
		context.Data["YourTeams"] = yourteams
	}

	context.Data["Teams"] = teams
	context.Render("event-dashboard")
}

/*
func voteTeam(rw http.ResponseWriter, r *http.Request) {
	Page(rw, r, "Vote for XYZ",
		html.Form().Child(
			fieldset(
				legend("Entry"),

				field("Name"),
				field("Team"),

				textarea("Guide", ""),

				field("Windows"),
				field("Linux"),
				field("Mac"),
				field("Web"),
			),

			fieldset(
				legend("Aspects"),
				aspect("Theme"),
				aspect("Enjoyment"),
				aspect("Aesthetics"),
				aspect("Innovation"),
				aspect("Bonus"),
				aspect("Total"),
			),

			html.Submit("Create"),
		),
	)
}
func fieldset(rs ...html.Renderer) *html.Node {
	return html.Tag("fieldset", "", rs...)
}

func createTeam(rw http.ResponseWriter, r *http.Request) {
	Page(rw, r, "New Entry",
		html.Form().Child(
			field("Event"),
			field("Name"),
			field("Members"),
			html.Submit("Create"),
		),
	)
}

func legend(text string) *html.Node {
	return html.Tag("legend", "", html.Text{text})
}

func field(label string) *html.Node {
	return html.Div("field",
		html.Label(label, label),
		html.Input(label, "text"),
	)
}

func aspect(label string) *html.Node {
	return html.Div("field",
		html.Label(label, label),
		html.Input(label, "range").
			Attr("min", "1").
			Attr("max", "5").
			Attr("step", "0.01"),
	)
}

func datetime(label string) *html.Node {
	return html.Div("field",
		html.Label(label, label),
		html.Input(label, "datetime-local"),
	)
}

func textarea(label, defaultValue string) *html.Node {
	return html.Div("field",
		html.Label(label, label),
		html.Textarea(label).Text(defaultValue),
	)
}
*/
