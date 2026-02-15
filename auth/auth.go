package auth

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
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
	Log *slog.Logger

	Development struct {
		Enabled  bool
		Sessions sessions.Store
	}

	OAuth2         *oauth2.Config
	Sessions       sessions.Store
	Domain         string
	LoginFailed    string
	LoginCompleted string
}

// NewService returns new Service for the given domain.
func NewService(log *slog.Logger, domain string, oauthConfig *oauth2.Config, sess sessions.Store) *Service {
	service := &Service{}
	service.Log = log

	service.Development.Enabled = false
	if service.Development.Enabled {
		cookieStore := sessions.NewCookieStore([]byte("DEVELOPMENT"))
		cookieStore.Options = &sessions.Options{
			Path:     "/",
			HttpOnly: true,
		}
		service.Development.Sessions = cookieStore
	}

	service.OAuth2 = oauthConfig
	service.Sessions = sess
	service.Domain = domain
	service.LoginFailed = "/"
	service.LoginCompleted = "/"

	return service
}

// Register registers handlers for /auth/*.
func (service *Service) Register(router *mux.Router) {
	router.HandleFunc("/auth/login", service.Login)
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
	return []Link{{"Google", "/auth/login"}}
}

// CurrentCredentials returns credentials associated with the request.
func (service *Service) CurrentCredentials(r *http.Request) *Credentials {
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

	sess, err := service.Sessions.Get(r, "jamvote")
	if err != nil {
		return nil
	}

	provider, _ := sess.Values["auth_provider"].(string)
	id, _ := sess.Values["auth_id"].(string)
	if provider == "" || id == "" {
		return nil
	}

	email, _ := sess.Values["auth_email"].(string)
	name, _ := sess.Values["auth_name"].(string)

	return &Credentials{
		Provider: provider,
		ID:       id,
		Email:    email,
		Name:     name,
		Admin:    false,
	}
}

// Login initiates the OAuth2 flow.
func (service *Service) Login(w http.ResponseWriter, r *http.Request) {
	if service.OAuth2.ClientID == "" {
		service.Log.Warn("OAuth2 not configured: missing client ID")
		http.Redirect(w, r, service.LoginFailed, http.StatusSeeOther)
		return
	}

	state, err := generateState()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	sess, _ := service.Sessions.Get(r, "jamvote")
	sess.Values["oauth_state"] = state
	if err := sess.Save(r, w); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	url := service.OAuth2.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

// idTokenClaims represents the relevant claims from a Google ID token.
type idTokenClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Callback is called after a login event.
func (service *Service) Callback(w http.ResponseWriter, r *http.Request) {
	sess, _ := service.Sessions.Get(r, "jamvote")

	// Verify CSRF state.
	expectedState, _ := sess.Values["oauth_state"].(string)
	delete(sess.Values, "oauth_state")
	if expectedState == "" || r.FormValue("state") != expectedState {
		http.Redirect(w, r, service.LoginFailed, http.StatusSeeOther)
		return
	}

	// Exchange code for token.
	token, err := service.OAuth2.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		service.Log.Error("OAuth2 exchange error", "error", err)
		http.Redirect(w, r, service.LoginFailed, http.StatusSeeOther)
		return
	}

	// Parse the id_token from the token response.
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		service.Log.Error("no id_token in OAuth2 response")
		http.Redirect(w, r, service.LoginFailed, http.StatusSeeOther)
		return
	}

	claims, err := parseIDToken(idToken)
	if err != nil {
		service.Log.Error("failed to parse id_token", "error", err)
		http.Redirect(w, r, service.LoginFailed, http.StatusSeeOther)
		return
	}

	// Store credentials in session.
	sess.Values["auth_provider"] = "google"
	sess.Values["auth_id"] = claims.Sub
	sess.Values["auth_email"] = claims.Email
	sess.Values["auth_name"] = claims.Name
	if err := sess.Save(r, w); err != nil {
		service.Log.Error("failed to save session", "error", err)
	}

	http.Redirect(w, r, service.LoginCompleted, http.StatusSeeOther)
}

// Logout is called when a user wants to log out.
func (service *Service) Logout(w http.ResponseWriter, r *http.Request) {
	if service.Development.Enabled {
		sess, _ := service.Development.Sessions.New(r, DefaultDevelopmentSession)
		sess.Values["User"] = ""
		if err := sess.Save(r, w); err != nil {
			service.Log.Error("failed to save session on logout", "error", err)
		}
	}

	sess, _ := service.Sessions.Get(r, "jamvote")
	delete(sess.Values, "auth_provider")
	delete(sess.Values, "auth_id")
	delete(sess.Values, "auth_email")
	delete(sess.Values, "auth_name")
	if err := sess.Save(r, w); err != nil {
		service.Log.Error("failed to clear session on logout", "error", err)
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
		service.Log.Error("unable to create development session", "error", err)
		http.Error(w, "unable to create development session", http.StatusInternalServerError)
		return
	}

	username := strings.TrimSpace(r.FormValue("name"))

	if username != "" {
		sess.Values["User"] = username
		err := sess.Save(r, w)
		if err != nil {
			service.Log.Error("unable to save development session", "error", err)
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

// generateState creates a random state string for CSRF protection.
func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// parseIDToken extracts claims from a JWT id_token by base64-decoding the payload segment.
func parseIDToken(idToken string) (*idTokenClaims, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) < 2 {
		return nil, errors.New("invalid id_token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	claims := &idTokenClaims{}
	if err := json.Unmarshal(payload, claims); err != nil {
		return nil, err
	}
	return claims, nil
}
