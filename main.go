package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/adinfinit/rater/html"
	"github.com/gorilla/pat"
)

var (
	listen = flag.String("listen", ":8080", "listen on address")
)

func main() {
	router := pat.New()
	router.Get("/user", user)
	router.Get("/", home)

	fmt.Printf("Listening on %q\n", *listen)
	err := http.ListenAndServe(*listen, router)
	if err != nil {
		log.Fatal(err)
	}
}

func Render(w http.ResponseWriter, r func(w *html.Writer)) {
	writer := html.NewWriter()
	r(writer)

	w.Write(writer.Bytes())
}

func home(rw http.ResponseWriter, r *http.Request) {
	Render(rw, func(w *html.Writer) {
		Page(w, "Rater",
			html.P("Lorem ipsum dolor sit amet, consectetur adipisicing elit. Reiciendis possimus quod repellendus hic consequatur aliquam unde velit harum, quae magnam dolorem alias odio, excepturi culpa est. Voluptates repellendus nihil quisquam!"),
		)
	})
}

func user(rw http.ResponseWriter, r *http.Request) {
	Render(rw, func(w *html.Writer) {
		Page(w, "User",
			html.P("Lorem ipsum dolor sit amet, consectetur adipisicing elit. Reiciendis possimus quod repellendus hic consequatur aliquam unde velit harum, quae magnam dolorem alias odio, excepturi culpa est. Voluptates repellendus nihil quisquam!"),
		)
	})
}

func Page(w *html.Writer, title string, content ...html.Renderer) {
	w.UnsafeWrite("<!DOCTYPE html>")
	w.Open("html")
	defer w.Close("html")

	w.Render(html.Head(
		html.Meta("charset", "utf-8"),
		html.Meta("http-equiv", "X-UA-Compatible").Attr("content", "IE=edge"),
		html.Meta("name", "google").Attr("value", "notranslate"),
		html.Title(title),
		html.Link("/static/reset.css"),
		html.Link("/static/main.css"),
	))

	w.Open("body")

	w.Render(
		html.Div("header",
			html.H1(title),
		),
		html.Div("content", content...),
	)

	w.Close("body")
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
