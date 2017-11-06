package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/adinfinit/rater/html"
)

var (
	listen = flag.String("listen", ":8080", "listen on address")
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", home)
	router.HandleFunc("/profile", profile)
	router.HandleFunc("/events", events)
	router.HandleFunc("/logout", logout)

	staticFiles := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	router.PathPrefix("/static/").Handler(staticFiles)

	fmt.Printf("Listening on %q\n", *listen)
	err := http.ListenAndServe(*listen, router)
	if err != nil {
		log.Fatal(err)
	}
}

func home(rw http.ResponseWriter, r *http.Request) {
	Page(rw, r, "Dashboard",
		html.P("Lorem ipsum dolor sit amet, consectetur adipisicing elit. Reiciendis possimus quod repellendus hic consequatur aliquam unde velit harum, quae magnam dolorem alias odio, excepturi culpa est. Voluptates repellendus nihil quisquam!"),
	)
}

func profile(rw http.ResponseWriter, r *http.Request) {
	Page(rw, r, "Profile",
		html.P("Lorem ipsum dolor sit amet, consectetur adipisicing elit. Reiciendis possimus quod repellendus hic consequatur aliquam unde velit harum, quae magnam dolorem alias odio, excepturi culpa est. Voluptates repellendus nihil quisquam!"),
	)
}

func events(rw http.ResponseWriter, r *http.Request) {
	Page(rw, r, "Events",
		html.P("Lorem ipsum dolor sit amet, consectetur adipisicing elit. Reiciendis possimus quod repellendus hic consequatur aliquam unde velit harum, quae magnam dolorem alias odio, excepturi culpa est. Voluptates repellendus nihil quisquam!"),
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
		html.Div("header",
			html.Div("title", html.Text{title}),
			Menu(r, [][2]string{
				{"/", "Dashboard"},
				{"/profile", "Profile"},
				{"/events", "Events"},
				{"/logout", "Sign out"},
			}),
		),
		html.Div("content", content...),
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

type UserID string
type TeamID string

type User struct {
	ID    UserID
	Name  string
	Teams []Team
}

type Team struct {
	ID   TeamID
	Name string

	Members []Member
}

type Member struct {
	ID   UserID
	Name string
}
