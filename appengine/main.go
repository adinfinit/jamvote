package main

import (
	"net/http"

	"google.golang.org/appengine"

	"github.com/gorilla/mux"

	"github.com/adinfinit/jamvote/auth"
	"github.com/adinfinit/jamvote/dashboard"
	"github.com/adinfinit/jamvote/user"
)

func main() {
	router := mux.NewRouter()

	auths := auth.NewService("http://localhost:8080")
	auths.LoginCompleted = "/user"
	auths.LoginFailed = "/user/login"
	auths.Register(router)

	users := &user.Server{auths}
	users.Register(router)

	dashboards := &dashboard.Server{users}
	dashboards.Register(router)

	// events := event.NewServer("LD40", "Ludum Dare 40")
	// events.Register(router)

	http.Handle("/", router)

	appengine.Main()
}
