package auth

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

// Credentials represents users authentication method.
type Credentials struct {
	Provider string
	ID       string
	Email    string
	Name     string
	Admin    bool
}

// DefaultDevelopmentSession is the default session name for development login.
const DefaultDevelopmentSession = "jamvote-development"

// Service implements authentication endpoints.
type Service struct {
	Development struct {
		Enabled  bool
		Sessions sessions.Store
	}

	Domain         string
	LoginFailed    string
	LoginCompleted string
}

// NewService returns new Service for the given domain.
func NewService(domain string) *Service {
	service := &Service{}

	service.Development.Enabled = false
	if service.Development.Enabled {
		cookieStore := sessions.NewCookieStore([]byte("DEVELOPMENT"))
		cookieStore.Options = &sessions.Options{
			Path:     "/",
			HttpOnly: true,
		}
		service.Development.Sessions = cookieStore
	}

	service.Domain = domain
	service.LoginFailed = "/"
	service.LoginCompleted = "/"

	return service
}

// Register registers handlers for /auth/*.
func (service *Service) Register(router *mux.Router) {
	router.HandleFunc("/auth/callback", service.Callback)
	router.HandleFunc("/auth/development-login", service.DevelopmentLogin)
	router.HandleFunc("/auth/logout", service.Logout)
}

// Link represents a single login URL for a provider.
type Link struct {
	Title string
	URL   string
}

// Links returns all available login URLs.
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

// CurrentCredentials returns credentials associated with the request.
func (service *Service) CurrentCredentials(c context.Context, r *http.Request) *Credentials {
	if service.Development.Enabled {
		sess, _ := service.Development.Sessions.New(r, DefaultDevelopmentSession)
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

// Callback is called after a login event.
func (service *Service) Callback(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	aeuser := user.Current(c)
	if aeuser != nil {
		http.Redirect(w, r, service.LoginCompleted, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, service.LoginFailed, http.StatusSeeOther)
	}
}

// Logout is called when a user wants to log out.
func (service *Service) Logout(w http.ResponseWriter, r *http.Request) {
	if service.Development.Enabled {
		sess, _ := service.Development.Sessions.New(r, DefaultDevelopmentSession)
		sess.Values["User"] = ""
		err := sess.Save(r, w)
		if err != nil {
			log.Println("Failed to logout:", err)
		}
	}

	c := appengine.NewContext(r)
	logout, err := user.LogoutURL(c, "/")
	if err == nil {
		http.Redirect(w, r, logout, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DevelopmentLogin is a way to login to a development server.
func (service *Service) DevelopmentLogin(w http.ResponseWriter, r *http.Request) {
	if !service.Development.Enabled {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	sess, err := service.Development.Sessions.New(r, DefaultDevelopmentSession)
	if err != nil {
		log.Println("Unable to create development session:", err)
		http.Error(w, "unable to create development session", http.StatusInternalServerError)
		return
	}

	username := strings.TrimSpace(r.FormValue("name"))

	if username != "" {
		sess.Values["User"] = username
		err := sess.Save(r, w)
		if err != nil {
			log.Println("Unable to save development session:", err)
		}
		http.Redirect(w, r, service.LoginCompleted, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, service.LoginFailed, http.StatusSeeOther)
	}
}

// developmentUserID dynamically creates a user ID.
func developmentUserID(username string) string {
	h := sha1.Sum([]byte(username))
	return hex.EncodeToString(h[:])
}
