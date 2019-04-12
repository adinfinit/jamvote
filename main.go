package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/adinfinit/jamvote/auth"
	"github.com/adinfinit/jamvote/datastoredb"
	"github.com/adinfinit/jamvote/event"
	"github.com/adinfinit/jamvote/profile"
	"github.com/adinfinit/jamvote/site"
	"github.com/adinfinit/jamvote/user"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

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

	users := &user.Server{
		Site: sites,
		DB:   db,
		Auth: auths,
	}
	users.Register(router)

	events := &event.Server{
		Site:  sites,
		DB:    db,
		Users: users,
	}
	events.Register(router)

	profiles := &profile.Server{
		Site:   sites,
		Events: db,
		Users:  users,
	}
	profiles.Register(router)

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
