package main

import (
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
	projectid := os.Getenv("GOOGLE_CLOUD_PROJECT")
	db := datastoredb.OpenDB(projectid)
	log.Printf("Opened database with project-id:%q", projectid)

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
