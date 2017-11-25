package event

import (
	"net/http"

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
	router.HandleFunc("/create-event", event.Users.Handler(event.CreateEvent))
	router.HandleFunc("/event/{eventid}", event.Handler(event.Dashboard))
	router.HandleFunc("/event/{eventid}/create-team", event.Handler(event.CreateTeam))
	router.HandleFunc("/event/{eventid}/progress", event.Handler(event.Dashboard))
	router.HandleFunc("/event/{eventid}/closing", event.Handler(event.Dashboard))
	router.HandleFunc("/event/{eventid}/summary", event.Handler(event.Dashboard))
	router.HandleFunc("/event/{eventid}/team/{teamid}", event.Handler(event.Team))
	router.HandleFunc("/event/{eventid}/team/{teamid}/edit", event.Handler(event.EditTeam))
	router.HandleFunc("/event/{eventid}/vote/{teamid}", event.Handler(event.Dashboard))
}

func (event *Server) Dashboard(context *Context) {
	context.Render("event-dashboard")
}

func (event *Server) Team(context *Context) {
	teamid, ok := context.IntParam("teamid")
	if !ok {
		context.Error("Team ID missing", http.StatusBadRequest)
		return
	}

	context.Data["Team"] = Team{
		ID: TeamID(teamid),
	}

	context.Render("event-team")
}

func (event *Server) EditTeam(context *Context) {
	teamid, ok := context.IntParam("teamid")
	if !ok {
		context.Error("Team ID missing", http.StatusBadRequest)
		return
	}

	context.Data["Team"] = Team{
		ID: TeamID(teamid),
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
