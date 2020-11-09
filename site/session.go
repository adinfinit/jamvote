package site

import (
	"crypto/rand"
	"log"
	"os"

	"github.com/gorilla/sessions"
)

// DefaultSessionName is the name of default session store.
const DefaultSessionName = "jamvote"

// SessionStore keeps all the alive sessions and cookies.
var SessionStore sessions.Store

func init() {
	secret := []byte(os.Getenv("COOKIESTORE_SECRET"))
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
	SessionStore = cookieStore
}
