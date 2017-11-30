package site

import (
	"context"
	"fmt"
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

var aptLocation *time.Location

func init() {
	var err error
	aptLocation, err = time.LoadLocation("Europe/Tallinn")
	if err != nil {
		panic(err)
	}
}

type Server struct {
	Development bool
}

type Context struct {
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
	if flashes := sess.Flashes(); len(flashes) > 0 {
		data["Flashes"] = flashes
	}
	data["Development"] = server.Development

	return &Context{
		Request:  r,
		Response: w,
		Data:     data,
		Session:  sess,
		Context:  appengine.NewContext(r),
	}
}

func (server *Server) Handler(fn func(*Context)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(server.Context(w, r))
	})
}

func (context *Context) Flash(message ...string) {
	for _, m := range message {
		context.Session.AddFlash(m)
	}
}

func (context *Context) FlashNow(message ...string) {
	if v, ok := context.Data["Flashes"]; ok {
		if list, ok := v.([]string); ok {
			list = append(list, message...)
			context.Data["Flashes"] = list
			return
		}
	}

	context.Data["Flashes"] = message
}

func (context *Context) Redirect(url string, status int) {
	context.Session.Save(context.Request, context.Response)
	http.Redirect(context.Response, context.Request, url, status)
}

func (context *Context) Error(text string, status int) {
	context.Session.Save(context.Request, context.Response)
	http.Error(context.Response, text, status)
}

func toFloat(a interface{}) float64 {
	switch v := a.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	}
	return 0
}

func (context *Context) Render(name string) {
	context.Session.Save(context.Request, context.Response)

	// TODO: cache template parsing
	t := template.New("")
	t = t.Funcs(template.FuncMap{
		"Data": func() interface{} { return nil },
		"formatDateTime": func(t *time.Time) string {
			// TODO: use event location
			return t.In(aptLocation).Format("2006-01-02 15:04:05 MST")
		},
		"paragraphs": func(s string) []string {
			s = strings.Replace(s, "\r", "", -1)
			return strings.Split(s, "\n\n")
		},

		"add": func(a, b interface{}) float64 {
			return toFloat(a) + toFloat(b)
		},
		"sub": func(a, b interface{}) float64 {
			return toFloat(a) - toFloat(b)
		},
		"mul": func(a, b interface{}) float64 {
			return toFloat(a) * toFloat(b)
		},
		"div": func(a, b interface{}) float64 {
			return toFloat(a) / toFloat(b)
		},
	})

	t, err := t.ParseGlob("templates/**/*.html")
	if err != nil {
		http.Error(context.Response, fmt.Sprintf("Template error: %q", err), http.StatusInternalServerError)
		return
	}

	t = t.Funcs(template.FuncMap{"Data": func() interface{} { return context.Data }})

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
