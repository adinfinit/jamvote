package site

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

var aptLocation *time.Location

func init() {
	var err error
	aptLocation, err = time.LoadLocation("Europe/Tallinn")
	if err != nil {
		panic(err)
	}
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
	context.saveSession()

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
