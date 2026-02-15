package main

import (
	"crypto/rand"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"google.golang.org/appengine/v2"

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

	sessionsStore := newCookieSessionStore(os.Getenv("COOKIESTORE_SECRET"))

	sites, err := site.NewServer(sessionsStore, "./static", "templates/**/*.html")
	if err != nil {
		log.Fatal(err)
	}
	sites.Register(router)

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
		log.Printf("Listening on port %s", port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		appengine.Main()
	}
}

func newCookieSessionStore(secretString string) sessions.Store {
	secret := []byte(secretString)
	if len(secret) == 0 {
		log.Println("Cookie Secret missing")
		var code [64]byte
		_, err := rand.Read(code[:])
		if err != nil {
			log.Println(err)
			secret = code[:]
		}
	}

	cookieStore := sessions.NewCookieStore(secret)
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
	}

	return cookieStore
}
