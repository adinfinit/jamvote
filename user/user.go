package user

import (
	"strconv"

	"github.com/adinfinit/jamvote/auth"
)

type Repo interface {
	ByCredentials(cred *auth.Credentials) (*User, error)
	ByID(id ID) (*User, error)
	List() ([]*User, error)

	Create(cred *auth.Credentials, user *User) (ID, error)
	Update(user *User) error
}

type User struct {
	ID    ID `datastore:"-"`
	Name  string
	Admin bool `datastore:",noindex"`

	Facebook string `datastore:",noindex"`
	Github   string `datastore:",noindex"`

	NewUser bool `datastore:"-"`
}

func (user *User) IsAdmin() bool {
	return user != nil && user.Admin
}

func (user *User) Equals(b *User) bool {
	return user.ID == b.ID
}

type ID int64

func (id ID) String() string { return strconv.Itoa(int(id)) }
