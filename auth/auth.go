package auth

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"

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

	Configs []*oauth2.Config
}

func NewService(domain string) *Service {
	service := &Service{}

	service.Domain = domain
	service.LoginFailed = "/"
	service.LoginCompleted = "/"

	service.AddDefaultProviders()

	return service
}

func (service *Service) Register(router *mux.Router) {
	router.HandleFunc("/auth/callback", service.Callback)
	router.HandleFunc("/auth/logout", service.Logout)
}

func (service *Service) add(name string, endpoint oauth2.Endpoint) {
	uname := strings.ToUpper(name)
	clientid := os.Getenv(uname + "_ID")
	secret := os.Getenv(uname + "_SECRET")

	if clientid != "" && secret != "" {
		service.Configs = append(service.Configs, &oauth2.Config{
			ClientID:     clientid,
			ClientSecret: secret,
			Scopes:       []string{"user"},
			Endpoint:     endpoint,
		})
	}
}

func (service *Service) AddDefaultProviders() {
	service.add("facebook", facebook.Endpoint)
	service.add("github", github.Endpoint)
}

type LoginLink struct {
	Title string
	URL   string
}

var providers = map[string]string{
	"Google":   "",
	"Github":   "github.com",
	"Facebook": "facebook.com",
}

func (service *Service) Logins(r *http.Request) []LoginLink {
	infos := []LoginLink{}

	c := appengine.NewContext(r)
	for _, provider := range []string{"Google", "Github", "Facebook"} {
		identity := providers[provider]
		loginurl, err := user.LoginURLFederated(c, "/auth/callback", identity)
		if err != nil {
			log.Println(err)
			continue
		}
		infos = append(infos, LoginLink{provider, loginurl})
	}

	return infos
}

func (service *Service) CurrentCredentials(w http.ResponseWriter, r *http.Request) *Credentials {
	c := appengine.NewContext(r)

	feduser, err := user.CurrentOAuth(c)
	if feduser != nil && err == nil {
		return &Credentials{
			Provider: "appengine",
			ID:       feduser.ID,
			Name:     feduser.Email,
			Email:    feduser.Email,
		}
	}

	return nil
}

func (service *Service) L(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	aeuser := user.Current(c)
	if aeuser != nil {
		http.Redirect(w, r, service.LoginCompleted, http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, service.LoginFailed, http.StatusTemporaryRedirect)
	}
}

func (service *Service) Callback(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	feduser, err := user.CurrentOAuth(c)
	if feduser != nil && err == nil {
		http.Redirect(w, r, service.LoginCompleted, http.StatusTemporaryRedirect)
	} else {
		log.Println(err)
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
