package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

type Credentials struct {
	Provider string
	ID       string
	Email    string
	Name     string
}

var (
	ErrNotLoggedIn = errors.New("not logged in")
)

type Service struct {
	Domain         string
	LoginFailed    string
	LoginCompleted string
}

func NewService(domain string) *Service {
	service := &Service{}

	service.Domain = domain
	service.LoginFailed = "/"
	service.LoginCompleted = "/"

	return service
}

func (service *Service) Register(router *mux.Router) {
	router.HandleFunc("/auth/callback", service.Callback)
	router.HandleFunc("/auth/logout", service.Logout)
}

type Link struct {
	Title string
	URL   string
}

func (service *Service) Links(r *http.Request) []Link {
	infos := []Link{}

	c := appengine.NewContext(r)
	loginurl, err := user.LoginURL(c, "/auth/callback")
	if err != nil {
		log.Println(err)
		return infos
	}
	infos = append(infos, Link{"Google", loginurl})

	return infos
}

func (service *Service) CurrentCredentials(r *http.Request) *Credentials {
	c := appengine.NewContext(r)

	aeuser := user.Current(c)
	if aeuser != nil {
		name := aeuser.Email
		if p := strings.Index(name, "@"); p >= 0 {
			name = name[:p]
		}

		return &Credentials{
			Provider: "appengine",
			ID:       aeuser.ID,
			Name:     name,
			Email:    aeuser.Email,
		}
	}

	return nil
}

func (service *Service) Callback(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	aeuser := user.Current(c)
	if aeuser != nil {
		http.Redirect(w, r, service.LoginCompleted, http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, service.LoginFailed, http.StatusTemporaryRedirect)
	}
}

func (service *Service) Logout(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	logout, err := user.LogoutURL(c, "/")
	if err == nil {
		http.Redirect(w, r, logout, http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
