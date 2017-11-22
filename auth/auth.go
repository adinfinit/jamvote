package auth

import (
	"errors"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gplus"
)

type User struct {
	Provider string
	Email    string
	Name     string
	UserID   string
}

var (
	ErrNotLoggedIn = errors.New("not logged in")
)

type Service struct {
	LoginFailed    string
	LoginCompleted string
}

func NewService() *Service {
	service := &Service{}

	service.LoginFailed = "/"
	service.LoginCompleted = "/"

	service.AddDefaultProviders()

	return service
}

func (service *Service) Register(router *mux.Router) {
	router.HandleFunc("/auth/{provider}", service.Login)
	router.HandleFunc("/auth/appengine/callback", service.LoginCallbackAppengine)
	router.HandleFunc("/auth/{provider}/callback", service.LoginCallback)
	router.HandleFunc("/auth/logout", service.Logout)
}

func (service *Service) add(name string, mk func(key, secret, url string) goth.Provider) {
	uname := strings.ToUpper(name)
	key := os.Getenv(uname + "_KEY")
	secret := os.Getenv(uname + "_SECRET")

	if key != "" && secret != "" {
		url := path.Join("/auth", name, "callback")
		provider := mk(key, secret, url)
		goth.UseProviders(provider)
	}
}

func (service *Service) AddDefaultProviders() {
	service.add("gplus", func(key, secret, url string) goth.Provider { return gplus.New(key, secret, url) })
	service.add("facebook", func(key, secret, url string) goth.Provider { return facebook.New(key, secret, url) })
	service.add("github", func(key, secret, url string) goth.Provider { return github.New(key, secret, url) })
}

type LoginLink struct {
	Title string
	URL   string
}

func (service *Service) Logins(r *http.Request) []LoginLink {
	infos := []LoginLink{}

	for name, _ := range goth.GetProviders() {
		info := LoginLink{}
		info.Title = name
		info.URL = path.Join("/auth", name)
		infos = append(infos, info)
	}

	c := appengine.NewContext(r)
	googlelogin, err := user.LoginURL(c, "/auth/appengine/callback")
	if err == nil {
		infos = append(infos, LoginLink{"Google", googlelogin})
	}

	sort.Slice(infos, func(i, k int) bool {
		return infos[i].Title < infos[k].Title
	})

	return infos
}

func (service *Service) User(r *http.Request) (*User, error) {
	return nil, nil
}

func (service *Service) Login(w http.ResponseWriter, r *http.Request) {
	if _, err := gothic.CompleteUserAuth(w, r); err == nil {
		// already logged in
		service.LoginCallback(w, r)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (service *Service) LoginCallbackAppengine(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	userinfo := user.Current(c)
	if userinfo != nil {
		http.Redirect(w, r, service.LoginCompleted, http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, service.LoginFailed, http.StatusTemporaryRedirect)
	}
}

func (service *Service) LoginCallback(w http.ResponseWriter, r *http.Request) {
	_, err := gothic.CompleteUserAuth(w, r)
	if err == nil {
		http.Redirect(w, r, service.LoginCompleted, http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, service.LoginFailed, http.StatusTemporaryRedirect)
	}
}

func (service *Service) Logout(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)

	c := appengine.NewContext(r)
	userinfo := user.Current(c)
	if userinfo != nil {
		logouturl, err := user.LogoutURL(c, "/")
		if err == nil {
			http.Redirect(w, r, logouturl, http.StatusTemporaryRedirect)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
