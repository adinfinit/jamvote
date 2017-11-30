package auth

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

type Credentials struct {
	Provider string
	ID       string
	Email    string
	Name     string
	Admin    bool
}

var (
	ErrNotLoggedIn = errors.New("not logged in")
)

type Service struct {
	Development bool

	Domain         string
	LoginFailed    string
	LoginCompleted string
}

func NewService(domain string) *Service {
	service := &Service{}

	service.Development = false

	service.Domain = domain
	service.LoginFailed = "/"
	service.LoginCompleted = "/"

	return service
}

func (service *Service) Register(router *mux.Router) {
	router.HandleFunc("/auth/callback", service.Callback)
	router.HandleFunc("/auth/development-login", service.DevelopmentLogin)
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

func (service *Service) CurrentCredentials(c context.Context, r *http.Request) *Credentials {
	if service.Development {
		sess, _ := developmentSessionStore.New(r, developmentSession)
		if val, ok := sess.Values["User"]; ok {
			if username, ok := val.(string); ok && username != "" {
				return &Credentials{
					Provider: "development",
					ID:       developmentUserID(username),
					Name:     username,
					Email:    username,
					Admin:    username == "Admin",
				}
			}
		}
	}

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
			Admin:    aeuser.Admin,
		}
	}

	return nil
}

func (service *Service) Callback(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	aeuser := user.Current(c)
	if aeuser != nil {
		http.Redirect(w, r, service.LoginCompleted, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, service.LoginFailed, http.StatusSeeOther)
	}
}

func (service *Service) Logout(w http.ResponseWriter, r *http.Request) {
	if service.Development {
		sess, _ := developmentSessionStore.New(r, developmentSession)
		sess.Values["User"] = ""
		sess.Save(r, w)
	}

	c := appengine.NewContext(r)
	logout, err := user.LogoutURL(c, "/")
	if err == nil {
		http.Redirect(w, r, logout, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (service *Service) DevelopmentLogin(w http.ResponseWriter, r *http.Request) {
	if !service.Development {
		return
	}

	sess, _ := developmentSessionStore.New(r, developmentSession)
	r.ParseForm()
	username := strings.TrimSpace(r.FormValue("name"))

	if username != "" {
		sess.Values["User"] = username
		sess.Save(r, w)
		http.Redirect(w, r, service.LoginCompleted, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, service.LoginFailed, http.StatusSeeOther)
	}
}

func developmentUserID(username string) string {
	h := sha1.Sum([]byte(username))
	return hex.EncodeToString(h[:])
}

const developmentSession = "jamvote-development"

var developmentSessionStore sessions.Store

func init() {
	cookieStore := sessions.NewCookieStore([]byte("DEVELOPMENT"))
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
	}
	developmentSessionStore = cookieStore
}
