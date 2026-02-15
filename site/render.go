package site

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"rsc.io/markdown"
)

// APTLocation is the location where times are handled.
var APTLocation *time.Location

func init() {
	var err error
	APTLocation, err = time.LoadLocation("Europe/Tallinn")
	if err != nil {
		panic(err)
	}
}

// templatesDir extracts the base directory from a glob pattern.
func templatesDir(glob string) string {
	dir := glob
	for _, ch := range []string{"*", "?", "["} {
		if i := strings.Index(dir, ch); i >= 0 {
			dir = dir[:i]
		}
	}
	return filepath.Dir(dir)
}

// renderMarkdown converts markdown text to HTML.
func renderMarkdown(s string) template.HTML {
	var p markdown.Parser
	doc := p.Parse(s)
	return template.HTML(markdown.ToHTML(doc))
}

// initTemplates loads all templates.
func (server *Server) initTemplates(glob string) error {
	server.TemplatesDir = templatesDir(glob)
	dir := server.TemplatesDir

	t := template.New("")
	t = t.Funcs(template.FuncMap{
		"Data": func() interface{} { return nil },
		"formatDateTime": func(t *time.Time) string {
			// TODO: use event location
			return t.In(APTLocation).Format("2006-01-02 15:04:05 MST")
		},
		"formatRFC": func(t time.Time) string {
			return t.In(APTLocation).Format(time.RFC3339)
		},
		"paragraphs": func(s string) []string {
			s = strings.Replace(s, "\r", "", -1)
			return strings.Split(s, "\n\n")
		},
		"ServerStartTime": func() string {
			return server.Start.Format("20060102T150405")
		},
		// hack to work around isZero
		"isValidTime": IsValidTime,

		"violinLeft":  violinLeft,
		"violinRight": violinRight,

		"averageViolinScore": func(min, max float64, xs interface{}) float64 {
			if x, ok := xs.(float64); ok {
				return 100 - 100*(x-min)/(max-min)
			}
			if xs, ok := xs.([]float64); ok {
				if len(xs) == 0 {
					return 50
				}
				t := 0.0
				for _, x := range xs {
					t += x
				}
				a := t / float64(len(xs))
				return 100 - 100*(a-min)/(max-min)
			}

			return -100
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
		"sequence1": func(count int) []int {
			xs := make([]int, count)
			for i := range xs {
				xs[i] = i + 1
			}
			return xs
		},
		"markdown": renderMarkdown,
		"markdownFile": func(name string) (template.HTML, error) {
			data, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				return "", err
			}
			return renderMarkdown(string(data)), nil
		},
	})

	t, err := t.ParseGlob(glob)
	server.Templates = t
	return err
}

// Render renders a template.
func (context *Context) Render(name string) {
	context.saveSession()
	t := context.Site.Templates.Funcs(template.FuncMap{"Data": func() interface{} { return context.Data }})
	if err := t.ExecuteTemplate(context.Response, name+".html", context.Data); err != nil {
		log.Println(err)
	}
}

// RenderMarkdown reads a markdown file and renders it wrapped with the markdown template.
func (context *Context) RenderMarkdown(name string) {
	data, err := os.ReadFile(filepath.Join(context.Site.TemplatesDir, name+".md"))
	if err != nil {
		log.Println(err)
		context.Error("page not found", http.StatusNotFound)
		return
	}

	context.Data["Content"] = renderMarkdown(string(data))
	context.Render("markdown")
}

// violinStep defines level of detail for violin plot.
const violinStep = 0.05

// violin calculates a violin plot for given scores.
func violin(min, max float64, scores []float64) []float64 {
	var values []float64
	maxvalue := 1.0

	for p := min; p <= max; p += violinStep {
		value := 0.0
		for _, score := range scores {
			value += cubicPulse(score, 0.5, p)
		}
		if maxvalue < value {
			maxvalue = value
		}
		values = append(values, value)
	}

	maxvalue *= 1.1
	for i := range values {
		values[i] /= maxvalue
	}

	return values
}

// violinLeft creates left side path of a violin plot.
func violinLeft(min, max float64, scores []float64) string {
	points := violin(min, max, scores)
	s := "50,100 "
	i := 0
	for p := min; p <= max; p += violinStep {
		s += fmt.Sprintf("%.0f,%.0f ", 50-points[i]*50, 100-100*(p-min)/(max-min))
		i++
	}
	s += "50,0"
	return s
}

// violinRight creates right side path of a violin plot.
func violinRight(min, max float64, scores []float64) string {
	points := violin(min, max, scores)
	s := "50,100 "
	i := 0
	for p := min; p <= max; p += violinStep {
		s += fmt.Sprintf("%.0f,%.0f ", 50+points[i]*50, 100-100*(p-min)/(max-min))
		i++
	}
	s += "50,0"
	return s
}

// cubicPulse calculates a smooth distance from a value.
func cubicPulse(center, radius, at float64) float64 {
	at = at - center
	if at < 0 {
		at = -at
	}
	if at > radius {
		return 0
	}
	at /= radius
	return 1 - at*at*(3-2*at)
}

// toFloat tries to convert a to a floating point value.
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
