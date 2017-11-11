package site

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Renderer struct {
	Glob string
}

func NewRenderer(glob string) *Renderer {
	r := &Renderer{}
	r.Glob = glob
	return r
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data interface{}) {
	t, err := template.ParseGlob("**/*.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %q", err), http.StatusInternalServerError)
		return
	}

	if err := t.ExecuteTemplate(w, name+".html", data); err != nil {
		log.Println(err)
	}
}
