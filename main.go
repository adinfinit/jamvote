package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/adinfinit/jamvote/about"
	"github.com/adinfinit/jamvote/auth"
	"github.com/adinfinit/jamvote/datastoredb"
	"github.com/adinfinit/jamvote/event"
	"github.com/adinfinit/jamvote/profile"
	"github.com/adinfinit/jamvote/site"
	"github.com/adinfinit/jamvote/user"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	ctx := context.Background()

	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	dsClient, err := datastore.NewClient(ctx, project)
	if err != nil {
		logger.Error("failed to create datastore client", "error", err)
		os.Exit(1)
	}
	defer dsClient.Close()

	db := &datastoredb.DB{Client: dsClient}

	router := mux.NewRouter()

	sessionsStore := newCookieSessionStore(logger, os.Getenv("COOKIESTORE_SECRET"))

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "http://localhost:8080"
	}

	oauthConfig := loadOAuthConfig(logger, domain+"/auth/callback")

	auths := auth.NewService(logger, domain, oauthConfig, sessionsStore)
	auths.LoginCompleted = "/user/logged-in"
	auths.LoginFailed = "/user/login"
	auths.Register(router)

	sites, err := site.NewServer(logger, sessionsStore, "./static", "templates/**/*.html")
	if err != nil {
		logger.Error("failed to create site server", "error", err)
		os.Exit(1)
	}
	sites.Register(router)

	users := &user.Server{
		Site: sites,
		DB:   db,
		Auth: auths,
	}
	users.Register(router)

	abouts := &about.Server{
		Site:  sites,
		Users: users,
	}
	abouts.Register(router)

	events := &event.Server{
		Site:  sites,
		DB:    db,
		Users: users,
	}
	events.Register(router)

	profiles := &profile.Server{
		Site:   sites,
		Events: db,
		Users:  users,
	}
	profiles.Register(router)

	http.Handle("/", router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Info("listening", "port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

// loadOAuthConfig loads the Google OAuth2 credentials JSON and returns an oauth2.Config.
// It checks GOOGLE_OAUTH_CREDENTIALS env var first (JSON content), then falls back
// to Secret Manager secret "GOOGLE_OAUTH_CREDENTIALS".
// The redirectURL overrides whatever is in the JSON file.
func loadOAuthConfig(logger *slog.Logger, redirectURL string) *oauth2.Config {
	jsonData := getSecretData(logger, "GOOGLE_OAUTH_CREDENTIALS")
	if len(jsonData) == 0 {
		logger.Warn("no Google OAuth credentials configured")
		return &oauth2.Config{RedirectURL: redirectURL}
	}

	cfg, err := google.ConfigFromJSON(jsonData, "openid", "email", "profile")
	if err != nil {
		logger.Error("failed to parse Google OAuth credentials JSON", "error", err)
		os.Exit(1)
	}
	cfg.RedirectURL = redirectURL
	return cfg
}

// getSecretData reads a secret value, checking the environment variable first,
// then falling back to Google Secret Manager.
func getSecretData(logger *slog.Logger, name string) []byte {
	if v := os.Getenv(name); v != "" {
		return []byte(v)
	}

	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if project == "" {
		return nil
	}

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		logger.Error("failed to create Secret Manager client", "error", err)
		return nil
	}
	defer client.Close()

	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", project, name),
	})
	if err != nil {
		logger.Error("failed to access secret", "secret", name, "error", err)
		return nil
	}

	return result.Payload.Data
}

func newCookieSessionStore(logger *slog.Logger, secretString string) sessions.Store {
	secret := []byte(secretString)
	if len(secret) == 0 {
		logger.Warn("cookie secret missing, generating random secret")
		var code [64]byte
		_, err := rand.Read(code[:])
		if err != nil {
			logger.Error("failed to generate random cookie secret", "error", err)
		}
		secret = code[:]
	}

	cookieStore := sessions.NewCookieStore(secret)
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
	}

	return cookieStore
}
