package user

import (
	"context"
	"encoding/gob"
	"errors"
	"strconv"

	"github.com/adinfinit/jamvote/auth"
)

// DB is interface for master database.
type DB interface {
	Users(context.Context) Repo
}

// Repo defines interaction with database.
type Repo interface {
	ByCredentials(cred *auth.Credentials) (*User, error)
	ByID(id UserID) (*User, error)
	List() ([]*User, error)

	Create(cred *auth.Credentials, user *User) (UserID, error)
	Update(user *User) error
}

// ErrNotExists is returned from Repo when a user does not exist.
var ErrNotExists = errors.New("user does not exist")

// UserID is a unique identifier for user.
type UserID int64

// String returns string representation of the id.
func (id UserID) String() string { return strconv.Itoa(int(id)) }

// User contains all relevant information about a user.
type User struct {
	ID    UserID `datastore:"-"`
	Name  string `datastore:",noindex"`
	Email string `datastore:",noindex"`
	Admin bool   `datastore:",noindex"`

	Facebook string `datastore:",noindex"`
	Github   string `datastore:",noindex"`

	NewUser bool `datastore:"-"`
}

func init() {
	gob.Register(&User{})
}

// IsAdmin returns whether a user is an administrator.
func (user *User) IsAdmin() bool {
	return user != nil && user.Admin
}

// HasEditor returns whether user can be edited by editor.
func (user *User) HasEditor(editor *User) bool {
	if editor.IsAdmin() {
		return true
	}
	return editor != nil && user.ID == editor.ID
}

// Equals returns whether b represent the same entity.
func (user *User) Equals(b *User) bool {
	return user.ID == b.ID
}
