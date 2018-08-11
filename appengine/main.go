package main

import (
	"log"
	"net/http"

	"google.golang.org/appengine"

	"github.com/gorilla/mux"

	"github.com/adinfinit/jamvote/auth"
	"github.com/adinfinit/jamvote/datastoredb"
	"github.com/adinfinit/jamvote/event"
	"github.com/adinfinit/jamvote/profile"
	"github.com/adinfinit/jamvote/site"
	"github.com/adinfinit/jamvote/user"
)

func main() {
	db := &datastoredb.DB{}

	router := mux.NewRouter()

	auths := auth.NewService("http://localhost:8080")
	auths.LoginCompleted = "/user/logged-in"
	auths.LoginFailed = "/user/login"
	auths.Register(router)

	sites, err := site.NewServer("templates/**/*.html")
	if err != nil {
		log.Fatal(err)
	}

	users := &user.Server{sites, db, auths}
	users.Register(router)

	events := &event.Server{sites, db, users}
	events.Register(router)

	profiles := &profile.Server{sites, db, users}
	profiles.Register(router)

	http.Handle("/", router)

	appengine.Main()
}
