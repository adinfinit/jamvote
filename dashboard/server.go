package dashboard

import (
	"github.com/adinfinit/jamvote/user"
	"github.com/gorilla/mux"
)

type Server struct {
	Users *user.Server
}

func (dashboard *Server) Register(router *mux.Router) {
	router.HandleFunc("/", dashboard.Users.Scoped(dashboard.Index))
}

func (dashboard *Server) Index(scope *user.Scope) {
	scope.Render("frontpage")
}
