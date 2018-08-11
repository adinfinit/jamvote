package user

import (
	"context"
	"encoding/gob"
	"errors"
	"strconv"

	"github.com/adinfinit/jamvote/auth"
)

type DB interface {
	Users(context.Context) Repo
}

type Repo interface {
	ByCredentials(cred *auth.Credentials) (*User, error)
	ByID(id UserID) (*User, error)
	List() ([]*User, error)

	Create(cred *auth.Credentials, user *User) (UserID, error)
	Update(user *User) error
}

var ErrNotExists = errors.New("info does not exist")

type UserID int64

func (id UserID) String() string { return strconv.Itoa(int(id)) }

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

func (user *User) IsAdmin() bool {
	return user != nil && user.Admin
}

func (user *User) HasEditor(editor *User) bool {
	if editor.IsAdmin() {
		return true
	}
	return editor != nil && user.ID == editor.ID
}

func (user *User) Equals(b *User) bool {
	return user.ID == b.ID
}
