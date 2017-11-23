package dashboard

import (
	"github.com/adinfinit/jamvote/user"
	"github.com/gorilla/mux"
)

type Server struct {
	Users *user.Server
}

func (dashboard *Server) Register(router *mux.Router) {
	router.HandleFunc("/", dashboard.Users.Handler(dashboard.Index))
}

func (dashboard *Server) Index(context *user.Context) {
	context.Render("frontpage")
}
