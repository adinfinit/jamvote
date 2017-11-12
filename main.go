package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/adinfinit/rater/html"
	"github.com/adinfinit/rater/site"
	"github.com/adinfinit/rater/user"
)

var (
	listen = flag.String("listen", ":8080", "listen on address")
)

func main() {
	router := mux.NewRouter()

	renderer := site.NewRenderer("**/*.html")

	mains := &site.Server{renderer}
	mains.Register(router)

	users := &user.Server{renderer}
	users.Register(router)

	router.HandleFunc("/ld40/", createTeam)
	router.HandleFunc("/ld40/team/create", createTeam)
	router.HandleFunc("/ld40/team/{teamid}", createTeam)
	router.HandleFunc("/ld40/vote/{teamid}", voteTeam)

	staticFiles := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	router.PathPrefix("/static/").Handler(staticFiles)

	fmt.Printf("Listening on %q\n", *listen)
	err := http.ListenAndServe(*listen, router)
	if err != nil {
		log.Fatal(err)
	}
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
	User(id user.ID) (user.User, error)
	Users() ([]user.User, error)

	Team(id TeamID) (Team, error)
	Teams() ([]Team, error)
}

type TeamID string

type Member struct {
	ID   user.ID
	Name string
}

type Event struct {
	Name string

	Create time.Time
	Start  time.Time
	Vote   time.Time
	Closed time.Time

	Organizers []user.ID
	Judges     []user.ID
	Teams      []TeamID
}

type Team struct {
	ID      TeamID
	Name    string
	Members []user.ID

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
	ID   user.ID
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
	clamp(&aspects.Bonus, 1, 5)
}

func (aspects *Aspects) Total() float64 {
	return (aspects.Theme +
		aspects.Enjoyment +
		aspects.Aesthetics +
		aspects.Innovation +
		aspects.Bonus*0.5) / (5*4 + 5*0.5)
}

func clamp(v *float64, min, max float64) {
	if *v < min {
		*v = min
	}
	if *v > max {
		*v = max
	}
}

/*

Theme
How well does it interpret the theme
1 Not even close
2 Resembling
3 Related
4 Spot on
5 Novel Interpretation

Enjoyment
How does the game generally feel
1 Boring
2 Not playing again
3 Nice
4 Didn't want to stop
5 Will play later.

Aesthetics
How well is the story, art and audio executed
1 None
2 Needs tweaks
3 Nice
4 Really good
5 Exceptional

Innovation
Something novel in the game
1 Seen it a lot
2 Interesting variation
3 Interesting approach
4 Never seen before
5 Exceptional

Bonus
Anything exceptionally special about * 0,5
1 Nothing special
2 Really liked *
3 Really loved *
4 Loved everything
5 <3

*/
