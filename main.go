package main

import (
	"fmt"
	"os"

	"github.com/adinfinit/rater/html"
)

func main() {
	w := html.NewWriter()
	HTML(w, "Hello World",
		html.H1("Hello"),
		html.P("Lorem ipsum dolor sit amet, consectetur adipisicing elit. Repellat, est, quas. Tempora veniam corrupti, alias, sunt doloremque magnam reiciendis ducimus commodi! Officiis commodi nihil repellendus repudiandae illum facilis, obcaecati quis!"),
	)
	os.Stdout.Write(w.Bytes())

	fmt.Println()
}

func HTML(w *html.Writer, title string, content ...html.Renderer) {
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
	w.Render(content...)
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
