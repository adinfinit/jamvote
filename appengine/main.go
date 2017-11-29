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
	auths.LoginCompleted = "/user/logged-in"
	auths.LoginFailed = "/user/login"
	auths.Register(router)

	sites := site.NewServer()
	sites.Global["Aspects"] = func() interface{} { return event.AspectsInfo }

	users := &user.Server{sites, auths}
	users.Register(router)

	events := &event.Server{sites, users}
	events.Register(router)

	http.Handle("/", router)

	appengine.Main()
}
