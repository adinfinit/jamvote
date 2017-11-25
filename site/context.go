package site

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	Data     map[string]interface{}

	context.Context
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Request:  r,
		Response: w,
		Data:     map[string]interface{}{},
		Context:  appengine.NewContext(r),
	}
}

func Handler(fn func(*Context)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(NewContext(w, r))
	})
}

func (context *Context) Redirect(url string, status int) {
	http.Redirect(context.Response, context.Request, url, status)
}

func (context *Context) Error(text string, status int) {
	http.Error(context.Response, text, status)
}

func (context *Context) Render(name string) {
	t := template.New("")
	t = t.Funcs(template.FuncMap{
		"formatDateTime": func(t *time.Time) string {
			return t.Format("2006-01-02 15:04:05 MST")
		},
	})

	t, err := t.ParseGlob("templates/**/*.html")
	if err != nil {
		http.Error(context.Response, fmt.Sprintf("Template error: %q", err), http.StatusInternalServerError)
		return
	}

	if err := t.ExecuteTemplate(context.Response, name+".html", context.Data); err != nil {
		log.Println(err)
	}
}

func (context *Context) StringParam(name string) (string, bool) {
	s, ok := mux.Vars(context.Request)[name]
	return s, ok
}

func (context *Context) IntParam(name string) (int64, bool) {
	s, ok := context.StringParam(name)
	if !ok {
		return 0, false
	}

	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return v, false
	}

	return v, true
}
