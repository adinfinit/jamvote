package main

import (
	"net/http"

	"google.golang.org/appengine"

	"github.com/gorilla/mux"

	"github.com/adinfinit/jamvote/auth"
	"github.com/adinfinit/jamvote/event"
	"github.com/adinfinit/jamvote/site"
	"github.com/adinfinit/jamvote/user"
)

func main() {
	router := mux.NewRouter()

	auths := auth.NewService("http://localhost:8080")
	auths.LoginCompleted = "/user"
	auths.LoginFailed = "/user/login"
	auths.Register(router)

	renderer := site.NewRenderer("templates/**/*.html")

	mains := &site.Server{renderer}
	mains.Register(router)

	users := &user.Server{auths, renderer}
	users.Register(router)

	events := event.NewServer("LD40", "Ludum Dare 40", renderer)
	events.Register(router)

	http.Handle("/", router)

	appengine.Main()
}
