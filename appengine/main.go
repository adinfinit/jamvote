package main

import (
	"flag"
	"net/http"

	"google.golang.org/appengine"

	"github.com/gorilla/mux"

	"github.com/adinfinit/rater/event"
	"github.com/adinfinit/rater/site"
	"github.com/adinfinit/rater/user"
)

var (
	listen = flag.String("listen", ":8080", "listen on address")
)

func main() {
	router := mux.NewRouter()

	renderer := site.NewRenderer("../**/*.html")

	mains := &site.Server{renderer}
	mains.Register(router)

	users := &user.Server{renderer}
	users.Register(router)

	events := event.NewServer("LD40", "Ludum Dare 40", renderer)
	events.Register(router)

	http.Handle("/", router)

	appengine.Main()
}
