package user

import "strconv"

type ID int64

func (id ID) String() string {
	return strconv.Itoa(int(id))
}

type User struct {
	ID       ID     `datastore:"-"`
	Name     string `datastore:",noindex"`
	Facebook string `datastore:",noindex"`
	Github   string `datastore:",noindex"`
	Admin    bool   `datastore:",noindex"`
}

func (user *User) Equals(b *User) bool {
	return user.ID == b.ID
}
