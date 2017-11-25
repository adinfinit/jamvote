package user

import (
	"strconv"

	"github.com/adinfinit/jamvote/auth"
)

type Repo interface {
	ByCredentials(cred *auth.Credentials) (*User, error)
	ByID(id ID) (*User, error)

	Create(cred *auth.Credentials, user *User) (ID, error)
	Update(user *User) error
}

type User struct {
	ID    ID     `datastore:"-"`
	Name  string `datastore:",noindex"`
	Admin bool   `datastore:",noindex"`

	Facebook string `datastore:",noindex"`
	Github   string `datastore:",noindex"`
}

func (user *User) Equals(b *User) bool {
	return user.ID == b.ID
}

type ID int64

func (id ID) String() string { return strconv.Itoa(int(id)) }
