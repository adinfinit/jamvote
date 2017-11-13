package event

import (
	"net/http"
	"path"

	"github.com/adinfinit/rater/site"
	"github.com/gorilla/mux"
)

type Server struct {
	Slug     string
	Title    string
	Renderer *site.Renderer
}

func NewServer(slug, title string, renderer *site.Renderer) *Server {
	return &Server{slug, title, renderer}
}

func (event *Server) Register(router *mux.Router) {
	prefix := path.Join("/", event.Slug)
	router.HandleFunc(prefix, event.Teams)
	router.HandleFunc(path.Join(prefix, "/create-team"), event.CreateTeam)
	router.HandleFunc(path.Join(prefix, "/progress"), event.Teams)
	router.HandleFunc(path.Join(prefix, "/closing"), event.Teams)
	router.HandleFunc(path.Join(prefix, "/summary"), event.Teams)
	router.HandleFunc(path.Join(prefix, "/{teamid}"), event.Team)
	router.HandleFunc(path.Join(prefix, "/vote/{teamid}"), event.Teams)
}

func (event *Server) Teams(w http.ResponseWriter, r *http.Request) {
	event.Renderer.Render(w, "event-teams", map[string]interface{}{
		"EventSlug":  event.Slug,
		"EventTitle": event.Title,
	})
}

func (event *Server) Team(w http.ResponseWriter, r *http.Request) {
	teamid := mux.Vars(r)["teamid"]

	event.Renderer.Render(w, "event-team", map[string]interface{}{
		"EventSlug":  event.Slug,
		"EventTitle": event.Title,

		"TeamID": teamid,
	})
}

func (event *Server) CreateTeam(w http.ResponseWriter, r *http.Request) {
	event.Renderer.Render(w, "create-team", map[string]interface{}{
		"EventSlug":  event.Slug,
		"EventTitle": event.Title,
	})
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
