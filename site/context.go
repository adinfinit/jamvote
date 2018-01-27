package site

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"google.golang.org/appengine"
)

type Server struct {
	Development bool
	Templates   *template.Template
}

func NewServer(templatesglob string) (*Server, error) {
	server := &Server{}
	return server, server.initTemplates(templatesglob)
}

type Context struct {
	Site *Server

	Request  *http.Request
	Response http.ResponseWriter
	Data     map[string]interface{}
	Session  *sessions.Session

	context.Context
}

func (server *Server) Context(w http.ResponseWriter, r *http.Request) *Context {
	sess, err := SessionStore.New(r, DefaultSessionName)
	if err != nil {
		log.Println("Failed to get session: ", err)
	}

	data := map[string]interface{}{}
	if flashes := sess.Flashes("_error"); len(flashes) > 0 {
		data["Errors"] = flashes
	}
	if flashes := sess.Flashes("_flash"); len(flashes) > 0 {
		data["Flashes"] = flashes
	}
	data["Development"] = server.Development

	return &Context{
		Site: server,

		Request:  r,
		Response: w,
		Data:     data,
		Session:  sess,
		Context:  appengine.NewContext(r),
	}
}

func (context *Context) saveSession() {
	context.Session.Save(context.Request, context.Response)
}

func (server *Server) Handler(fn func(*Context)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(server.Context(w, r))
	})
}

func (context *Context) flash(tag string, message ...string) {
	for _, m := range message {
		context.Session.AddFlash(m, tag)
	}
}

func (context *Context) flashNow(tag string, message ...string) {
	if v, ok := context.Data[tag]; ok {
		if list, ok := v.([]string); ok {
			list = append(list, message...)
			context.Data[tag] = list
			return
		}
	}

	context.Data[tag] = message
}

func (context *Context) FlashError(message ...string) {
	context.flash("_error", message...)
}
func (context *Context) FlashErrorNow(message ...string) {
	context.flashNow("Errors", message...)
}

func (context *Context) FlashMessage(message ...string) {
	context.flash("_flash", message...)
}
func (context *Context) FlashMessageNow(message ...string) {
	context.flashNow("Flashes", message...)
}

func (context *Context) Redirect(url string, status int) {
	context.saveSession()
	http.Redirect(context.Response, context.Request, url, status)
}

func (context *Context) Error(text string, status int) {
	context.saveSession()
	http.Error(context.Response, text, status)
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

func (context *Context) FormValue(name string) string {
	return strings.TrimSpace(context.Request.FormValue(name))
}
