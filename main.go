package jamvote

import (
	"flag"
	"fmt"
	"log"
	"net/http"

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

	renderer := site.NewRenderer("**/*.html")

	auths := auth.NewService()
	auths.Register(router)

	mains := &site.Server{renderer}
	mains.Register(router)

	users := &user.Server{renderer}
	users.Register(router)

	events := event.NewServer("LD40", "Ludum Dare 40", renderer)
	events.Register(router)

	staticFiles := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	router.PathPrefix("/static/").Handler(staticFiles)

	fmt.Printf("Listening on %q\n", *listen)

	err := http.ListenAndServe(*listen, router)
	if err != nil {
		log.Fatal(err)
	}
}
