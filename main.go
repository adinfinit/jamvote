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

	router.HandleFunc("/event/create", createEvent)
	router.HandleFunc("/event/{eventid}", createEvent)

	router.HandleFunc("/team/{eventid}/create", createTeam)
	router.HandleFunc("/team/{eventid}/{teamid}", createTeam)

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
	Page(rw, r, "New Event",
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

			fieldset(
				legend("Setup"),
				textarea("Aspects", "Aesthetics,1,5\nGraphics,1,5"),
				textarea("Fields", ""),
			),

			html.Submit("Create"),
		),
	)
}

func fieldset(rs ...html.Renderer) *html.Node {
	return html.Tag("fieldset", "", rs...)
}

func createTeam(rw http.ResponseWriter, r *http.Request) {
	Page(rw, r, "New Team",
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
				html.Div("title", html.Text{title}),
				Menu(r, [][2]string{
					{"/", "Events"},
					{"/user", "Profile"},
					{"/user/login", "Sign in"},
					{"/user/logout", "Sign out"},
				}),
			),
		),
		html.Div("content-outer",
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

	Event(id EventID) (Event, error)
	Events() ([]Event, error)
}

type UserID string
type TeamID string
type EventID string
type EntryID string

type User struct {
	ID    UserID
	Name  string
	Teams []TeamID
}

type Team struct {
	ID      TeamID
	Event   EventID
	Name    string
	Members []Member
}

type Member struct {
	ID   UserID
	Name string
}

type Event struct {
	ID   EventID
	Name string

	Create time.Time
	Start  time.Time
	Vote   time.Time
	Closed time.Time

	Fields  []Field
	Aspects []Aspect

	Organizers []UserID
	Judges     []UserID
	Teams      []TeamID

	Entries []EntryID
}

type Field struct {
	Name string
	Kind string
}

type Aspect struct {
	Name     string
	Min, Max float64
	Scale    float64
}

type Entry struct {
	ID    EntryID
	Event EventID

	Team    string
	Name    string
	Members []UserID

	Fields map[string]string
}

type Vote struct {
	ID      UserID
	Event   EventID
	Entry   EntryID
	Aspects map[string]float64

	Override bool
	Total    float64
}
