package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/adinfinit/rater/html"
)

var (
	listen = flag.String("listen", ":8080", "listen on address")
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", events)

	router.HandleFunc("/user/login", login)
	router.HandleFunc("/user/logout", logout)

	router.HandleFunc("/user", profile) // REDIRECT TO SELF
	router.HandleFunc("/user/{userid}", profile)

	router.HandleFunc("/team/create", createTeam)
	router.HandleFunc("/team/{teamid}", createTeam)
	router.HandleFunc("/vote/{teamid}", voteTeam)

	staticFiles := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	router.PathPrefix("/static/").Handler(staticFiles)

	fmt.Printf("Listening on %q\n", *listen)
	err := http.ListenAndServe(*listen, router)
	if err != nil {
		log.Fatal(err)
	}
}

func events(rw http.ResponseWriter, r *http.Request) {
	Page(rw, r, "Events",
		html.Section("events",
			html.H1("Ongoing"),
			eventLink("Ludum Dare 39", "Lorem", true),
		),
		html.Section("events events-small",
			html.H1("Completed"),
			eventLink("Ludum Dare 38", "Ipsum", false),
			eventLink("Ludum Dare 37", "Dolorem", false),
			eventLink("Ludum Dare 36", "Sigma", false),
			eventLink("Ludum Dare 35", "Delta", false),
			eventLink("Ludum Dare 34", "Phi", false),
		),
	)
}

func eventLink(title, theme string, ongoing bool) *html.Node {
	ongoingClass := ""
	if ongoing {
		ongoingClass = "ongoing"
	}
	return html.A("/event/123").Class("event-link").Class(ongoingClass).Child(
		html.Div("title").Text(title+" | "+theme),
		html.Div("countdown").Text("20:30:10"),
	)
}

func profile(rw http.ResponseWriter, r *http.Request) {
	Page(rw, r, "Profile",
		html.P("Lorem ipsum dolor sit amet, consectetur adipisicing elit. Debitis ipsa quidem itaque natus similique nemo voluptatum beatae doloremque, tempore blanditiis quod maiores quas tempora ad nesciunt officia accusamus atque veritatis!"),
	)
}

func createEvent(rw http.ResponseWriter, r *http.Request) {
	Page(rw, r, "Create Event",
		html.Form().Child(
			fieldset(
				legend("General"),
				field("Title"),
				field("Theme"),
			),

			fieldset(
				legend("Schedule"),
				datetime("Jam Start"),
				datetime("Voting Start"),
				datetime("Voting Closed"),
			),

			html.Submit("Create"),
		),
	)
}

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

func login(rw http.ResponseWriter, r *http.Request) {
	Page(rw, r, "Sign in",
		html.Div("logins",
			html.A("/user/login/google", html.Text{"Google"}),
			html.A("/user/login/facebook", html.Text{"Facebook"}),
		),
	)
}

func logout(rw http.ResponseWriter, r *http.Request) {
	http.Redirect(rw, r, "/", http.StatusFound)
}

func Page(rw http.ResponseWriter, r *http.Request, title string, content ...html.Renderer) {
	w := html.NewWriter()
	defer func() {
		rw.Write(w.Bytes())
	}()

	w.UnsafeWrite("<!DOCTYPE html>")
	w.Open("html")
	defer w.Close("html")

	w.Render(html.Head(
		html.Meta("charset", "utf-8"),
		html.Meta("http-equiv", "X-UA-Compatible").Attr("content", "IE=edge"),
		html.Meta("name", "google").Attr("value", "notranslate"),
		html.Title("Jam - "+title),
		html.Link("/static/main.css"),
	))

	w.Open("body")

	w.Render(
		html.Div("header-outer",
			html.Div("header",
				html.Div("title", html.Text{"Jamerator"}),
				Menu(r, [][2]string{
					{"/", "Events"},
					{"/user", "Profile"},
					{"/user/login", "Sign in"},
					{"/user/logout", "Sign out"},
				}),
			),
		),
		html.Div("content-outer",
			//html.Div("content", html.Text{title}),
			html.H1(title).Class("content"),
			html.Div("content", content...),
		),
	)

	w.Close("body")
}

func Menu(r *http.Request, items [][2]string) *html.Node {
	menu := html.Div("menu")
	for _, item := range items {
		link := html.A(item[0], html.Text{item[1]})
		if item[0] == r.URL.Path {
			link.Class("active")
		}
		menu.Child(link)
	}
	return menu
}

type Service interface {
	User(id UserID) (User, error)
	Users() ([]User, error)

	Team(id TeamID) (Team, error)
	Teams() ([]Team, error)
}

type UserID string
type TeamID string

type User struct {
	ID    UserID
	Name  string
	Teams []TeamID
}

type Member struct {
	ID   UserID
	Name string
}

type Event struct {
	Name string

	Create time.Time
	Start  time.Time
	Vote   time.Time
	Closed time.Time

	Organizers []UserID
	Judges     []UserID
	Teams      []TeamID
}

type Team struct {
	ID      TeamID
	Name    string
	Members []UserID

	Entry struct {
		Name         string
		Instructions string

		Link struct {
			Win string
			Mac string
			Web string
		}
	}
}

type Vote struct {
	ID   UserID
	Team TeamID

	Aspects  Aspects
	Override bool
	Total    float64
}

type Aspects struct {
	Theme      float64
	Enjoyment  float64
	Aesthetics float64
	Innovation float64
	Bonus      float64
}

func (aspects *Aspects) EnsureRange() {
	clamp(&aspects.Theme, 1, 5)
	clamp(&aspects.Enjoyment, 1, 5)
	clamp(&aspects.Aesthetics, 1, 5)
	clamp(&aspects.Innovation, 1, 5)
	clamp(&aspects.Bonus, 1, 3)
}

func (aspects *Aspects) Total() float64 {
	return (aspects.Theme +
		aspects.Enjoyment +
		aspects.Aesthetics +
		aspects.Innovation +
		aspects.Bonus) / (5*4 + 3)
}

func clamp(v *float64, min, max float64) {
	if *v < min {
		*v = min
	}
	if *v > max {
		*v = max
	}
}
