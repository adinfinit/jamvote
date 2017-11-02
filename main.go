package main

import (
	"fmt"
	"os"

	"github.com/adinfinit/rater/h"
)

func main() {
	w := h.NewWriter()
	HTML(w, "Hello World",
		h.H1("Hello"),
		h.P("Lorem ipsum dolor sit amet, consectetur adipisicing elit. Repellat, est, quas. Tempora veniam corrupti, alias, sunt doloremque magnam reiciendis ducimus commodi! Officiis commodi nihil repellendus repudiandae illum facilis, obcaecati quis!"),
	)
	os.Stdout.Write(w.Bytes())
	fmt.Println()
}

func HTML(w *h.Writer, title string, content ...h.Renderer) {
	w.UnsafeWrite("<!DOCTYPE html>")
	w.Open("html")
	defer w.Close("html")

	w.Render(h.Head(
		h.Meta("charset", "utf-8"),
		h.Meta("http-equiv", "X-UA-Compatible").Attr("content", "IE=edge"),
		h.Meta("name", "google").Attr("value", "notranslate"),
		h.Title(title),
		h.Link("/static/reset.css"),
		h.Link("/static/main.css"),
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
