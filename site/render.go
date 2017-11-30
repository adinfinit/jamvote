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

		"violinLeft":  violinLeft,
		"violinRight": violinRight,

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

const violinStep = 0.1

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
