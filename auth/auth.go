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

// Server implements authentication endpoints.
type Server struct {
	Development bool

	Domain         string
	LoginFailed    string
	LoginCompleted string
}

// NewServer returns new Server for the given domain.
func NewServer(domain string) *Server {
	service := &Server{}

	service.Development = false

	service.Domain = domain
	service.LoginFailed = "/"
	service.LoginCompleted = "/"

	return service
}

// Register registers handlers for /auth/*.
func (service *Server) Register(router *mux.Router) {
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
func (service *Server) Links(r *http.Request) []Link {
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
func (service *Server) CurrentCredentials(c context.Context, r *http.Request) *Credentials {
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

// Callback is called after a login event.
func (service *Server) Callback(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	aeuser := user.Current(c)
	if aeuser != nil {
		http.Redirect(w, r, service.LoginCompleted, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, service.LoginFailed, http.StatusSeeOther)
	}
}

// Logout is called when a user wants to log out.
func (service *Server) Logout(w http.ResponseWriter, r *http.Request) {
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

// DevelopmentLogin is a way to login to a development server.
func (service *Server) DevelopmentLogin(w http.ResponseWriter, r *http.Request) {
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

// developmentUserID dynamically creates a user ID.
func developmentUserID(username string) string {
	h := sha1.Sum([]byte(username))
	return hex.EncodeToString(h[:])
}

const developmentSession = "jamvote-development"

// developmentSessionStore is used for development logins.
var developmentSessionStore sessions.Store

func init() {
	cookieStore := sessions.NewCookieStore([]byte("DEVELOPMENT"))
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
	}
	developmentSessionStore = cookieStore
}
