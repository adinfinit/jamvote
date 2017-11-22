package site

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Renderer *Renderer
}

func (front *Server) Register(router *mux.Router) {
	router.HandleFunc("/", front.Frontpage)
}

func (front *Server) Frontpage(w http.ResponseWriter, r *http.Request) {
	front.Renderer.Render(w, "site-frontpage", map[string]interface{}{
		"UserID": "",
	})
}
