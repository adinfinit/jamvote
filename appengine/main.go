package main

import (
	"flag"
	"net/http"

	"google.golang.org/appengine"

	"github.com/gorilla/mux"

	"github.com/adinfinit/jamvote/auth"
	"github.com/adinfinit/jamvote/event"
	"github.com/adinfinit/jamvote/site"
	"github.com/adinfinit/jamvote/user"
)

var (
	listen = flag.String("listen", ":8080", "listen on address")
)

func main() {
	router := mux.NewRouter()

	auths := auth.NewService()
	auths.LoginCompleted = "/user"
	auths.LoginFailed = "/login"
	auths.Register(router)

	renderer := site.NewRenderer("../**/*.html")

	mains := &site.Server{renderer}
	mains.Register(router)

	users := &user.Server{auths, renderer}
	users.Register(router)

	events := event.NewServer("LD40", "Ludum Dare 40", renderer)
	events.Register(router)

	http.Handle("/", router)

	appengine.Main()
}
