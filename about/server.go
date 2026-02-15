package about

import (
	"github.com/adinfinit/jamvote/site"
	"github.com/adinfinit/jamvote/user"

	"github.com/gorilla/mux"
)

// Server implements profile views for users.
type Server struct {
	Site  *site.Server
	Users *user.Server
}

// Register registers user related endpoints.
func (server *Server) Register(router *mux.Router) {
	router.HandleFunc("/about", server.Handler(server.About))
	router.HandleFunc("/about/jamming", server.Handler(server.Jamming))
	router.HandleFunc("/about/scoring", server.Handler(server.Scoring))
	router.HandleFunc("/about/organizing", server.Handler(server.Organizing))
}

// About renders about page.
func (server *Server) About(context *Context) {
	context.RenderMarkdown("about/about")
}

// Jamming renders jamming page.
func (server *Server) Jamming(context *Context) {
	context.RenderMarkdown("about/about-jamming")
}

// Scoring renders scoring page.
func (server *Server) Scoring(context *Context) {
	context.RenderMarkdown("about/about-scoring")
}

// Organizing renders organizing page.
func (server *Server) Organizing(context *Context) {
	context.RenderMarkdown("about/about-organizing")
}
