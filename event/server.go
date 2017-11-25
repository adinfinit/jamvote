package event

import (
	"fmt"
	"net/http"
	"path"

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
	router.HandleFunc("/create-event", event.HandlerMaybe(event.CreateEvent))

	router.HandleFunc("/event/{eventid}", event.Handler(event.Dashboard))
	router.HandleFunc("/event/{eventid}/create-team", event.Handler(event.CreateTeam))
	router.HandleFunc("/event/{eventid}/voting", event.Handler(event.Dashboard))
	router.HandleFunc("/event/{eventid}/results", event.Handler(event.Dashboard))

	router.HandleFunc("/event/{eventid}/team/{teamid}", event.Handler(event.Team))
	router.HandleFunc("/event/{eventid}/team/{teamid}/edit", event.Handler(event.EditTeam))
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
		//TODO: flash
		context.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	context.Data["Events"] = events
	context.Render("event-list")
}

func (event *Server) Dashboard(context *Context) {
	context.Render("event-dashboard")
}

func (event *Server) Team(context *Context) {
	if context.Team == nil {
		context.Error("Team missing", http.StatusBadRequest)
		return
	}
	context.Render("event-team")
}

func (event *Server) EditTeam(context *Context) {
	if context.Team == nil {
		context.Error("Team missing", http.StatusBadRequest)
		return
	}

	context.Render("event-team")
}

func (event *Server) CreateTeam(context *Context) {
	if context.CurrentUser == nil {
		// TODO: add eventual return address
		context.Redirect("/user/login", http.StatusTemporaryRedirect)
		return
	}

	users, err := context.Users.List()
	if err != nil {
		context.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	context.Data["Users"] = users
	context.Render("event-create-team")
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
