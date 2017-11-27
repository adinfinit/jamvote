package site

import (
	"github.com/gorilla/sessions"
)

const DefaultSessionName = "jamvote"

var SessionStore sessions.Store

func init() {
	cookieStore := sessions.NewCookieStore([]byte("TODO TODO TODO TODO TODO"))
	cookieStore.Options = &sessions.Options{
		HttpOnly: true,
	}
	SessionStore = cookieStore
}
