package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
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
	ctx := context.Background()

	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	dsClient, err := datastore.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("Failed to create datastore client: %v", err)
	}
	defer dsClient.Close()

	db := &datastoredb.DB{Client: dsClient}

	router := mux.NewRouter()

	sessionsStore := newCookieSessionStore(os.Getenv("COOKIESTORE_SECRET"))

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "http://localhost:8080"
	}

	oauthConfig := loadOAuthConfig(domain + "/auth/callback")

	auths := auth.NewService(domain, oauthConfig, sessionsStore)
	auths.LoginCompleted = "/user/logged-in"
	auths.LoginFailed = "/user/login"
	auths.Register(router)

	sites, err := site.NewServer(sessionsStore, "./static", "templates/**/*.html")
	if err != nil {
		log.Fatal(err)
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
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// loadOAuthConfig loads the Google OAuth2 credentials JSON and returns an oauth2.Config.
// It checks GOOGLE_OAUTH_CREDENTIALS env var first (JSON content), then falls back
// to Secret Manager secret "GOOGLE_OAUTH_CREDENTIALS".
// The redirectURL overrides whatever is in the JSON file.
func loadOAuthConfig(redirectURL string) *oauth2.Config {
	jsonData := getSecretData("GOOGLE_OAUTH_CREDENTIALS")
	if len(jsonData) == 0 {
		log.Println("Warning: no Google OAuth credentials configured")
		return &oauth2.Config{RedirectURL: redirectURL}
	}

	cfg, err := google.ConfigFromJSON(jsonData, "openid", "email", "profile")
	if err != nil {
		log.Fatalf("Failed to parse Google OAuth credentials JSON: %v", err)
	}
	cfg.RedirectURL = redirectURL
	return cfg
}

// getSecretData reads a secret value, checking the environment variable first,
// then falling back to Google Secret Manager.
func getSecretData(name string) []byte {
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
		log.Printf("Failed to create Secret Manager client: %v", err)
		return nil
	}
	defer client.Close()

	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", project, name),
	})
	if err != nil {
		log.Printf("Failed to access secret %s: %v", name, err)
		return nil
	}

	return result.Payload.Data
}

func newCookieSessionStore(secretString string) sessions.Store {
	secret := []byte(secretString)
	if len(secret) == 0 {
		log.Println("Cookie Secret missing")
		var code [64]byte
		_, err := rand.Read(code[:])
		if err != nil {
			log.Println(err)
			secret = code[:]
		}
	}

	cookieStore := sessions.NewCookieStore(secret)
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
	}

	return cookieStore
}
