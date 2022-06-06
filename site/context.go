package site

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"google.golang.org/appengine"
)

// DefaultSessionName is the name of default session store.
const DefaultSessionName = "jamvote"

// Server implements the root request handler.
type Server struct {
	Development bool
	Start       time.Time
	Static      http.FileSystem
	Templates   *template.Template

	Sessions sessions.Store
}

// NewServer creates a new server with the specified templates.
func NewServer(sess sessions.Store, staticdir string, templatesglob string) (*Server, error) {
	server := &Server{}
	server.Sessions = sess
	server.Static = http.Dir(staticdir)
	server.Start = time.Now()
	return server, server.initTemplates(templatesglob)
}

func (server *Server) Register(router *mux.Router) {
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(server.Static)))
	router.Handle("/favicon.png", http.FileServer(server.Static))
	router.Handle("/robots.txt", http.FileServer(server.Static))
}

// Context contains all information about a single request.
type Context struct {
	Site *Server

	Request  *http.Request
	Response http.ResponseWriter
	Data     map[string]interface{}
	Session  *sessions.Session

	context.Context
}

// Context creates a new Context for a given request.
func (server *Server) Context(w http.ResponseWriter, r *http.Request) *Context {
	sess, err := server.Sessions.New(r, DefaultSessionName)
	if err != nil {
		log.Println("Failed to get session:", err)
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

// saveSession saves flashes for a given request.
func (context *Context) saveSession() {
	err := context.Session.Save(context.Request, context.Response)
	if err != nil {
		log.Println("Failed to save session:", err)
	}
}

// Handler wraps fn with automatic Context creation.
func (server *Server) Handler(fn func(*Context)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(server.Context(w, r))
	})
}

// flash adds a flashes with the specific tag.
func (context *Context) flash(tag string, message ...string) {
	for _, m := range message {
		context.Session.AddFlash(m, tag)
	}
}

// flashNow adds a flash to the current request.
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

// FlashError adds error messages.
func (context *Context) FlashError(message ...string) {
	context.flash("_error", message...)
}

// FlashErrorNow adds error messages for current request.
func (context *Context) FlashErrorNow(message ...string) {
	context.flashNow("Errors", message...)
}

// FlashMessage adds regular message.
func (context *Context) FlashMessage(message ...string) {
	context.flash("_flash", message...)
}

// FlashMessage adds regular message for current request.
func (context *Context) FlashMessageNow(message ...string) {
	context.flashNow("Flashes", message...)
}

// Redirect redirects to another url.
func (context *Context) Redirect(url string, status int) {
	context.saveSession()
	http.Redirect(context.Response, context.Request, url, status)
}

// Error responds with error code.
func (context *Context) Error(text string, status int) {
	context.saveSession()
	http.Error(context.Response, text, status)
}

// StringParam returns a string argument.
func (context *Context) StringParam(name string) (string, bool) {
	s, ok := mux.Vars(context.Request)[name]
	return s, ok
}

// IntParam returns an integer param.
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

// FormValue returns a form value.
func (context *Context) FormValue(name string) string {
	return strings.TrimSpace(context.Request.FormValue(name))
}
