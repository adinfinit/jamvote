package event

import (
	"net/http"
	"path"

	"github.com/adinfinit/jamvote/user"
	"github.com/gorilla/mux"
)

type Server struct {
	Slug  string
	Title string
	Users *user.Server
}

func NewServer(slug, title string, users *user.Server) *Server {
	return &Server{slug, title, users}
}

func (event *Server) Register(router *mux.Router) {
	prefix := path.Join("/", event.Slug)

	router.HandleFunc(prefix, event.Scoped(event.Teams))
	router.HandleFunc(path.Join(prefix, "/create-team"), event.Scoped(event.CreateTeam))
	router.HandleFunc(path.Join(prefix, "/progress"), event.Scoped(event.Teams))
	router.HandleFunc(path.Join(prefix, "/closing"), event.Scoped(event.Teams))
	router.HandleFunc(path.Join(prefix, "/summary"), event.Scoped(event.Teams))
	router.HandleFunc(path.Join(prefix, "/{teamid}"), event.Scoped(event.Team))
	router.HandleFunc(path.Join(prefix, "/vote/{teamid}"), event.Scoped(event.Teams))
}

func (event *Server) Teams(scope *Scope) {
	scope.Render("event-teams")
}

func (event *Server) Team(scope *Scope) {
	teamid, ok := scope.StringParam("teamid")
	if !ok {
		scope.Error("Team ID missing", http.StatusBadRequest)
		return
	}

	scope.Data["Team"] = Team{
		ID: TeamID(teamid),
	}

	scope.Render("event-team")
}

func (event *Server) CreateTeam(scope *Scope) {
	scope.Render("event-create-team")
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
