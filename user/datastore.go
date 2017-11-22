package user

import (
	"github.com/adinfinit/jamvote/auth"
	"google.golang.org/appengine/datastore"
)

type Datastore struct {
	Scope *Scope
}

type credentialMapping struct {
	UserKey  *datastore.Key
	Provider string `datastore:",noindex"`
	Email    string `datastore:",noindex"`
	Name     string `datastore:",noindex"`
}

func datastoreError(err error) error {
	if err == datastore.ErrNoSuchEntity {
		return ErrNotExists
	}
	return err
}

func (repo *Datastore) Create(cred *auth.Credentials, user *User) (ID, error) {
	incompletekey := datastore.NewIncompleteKey(repo.Scope, "User", nil)
	userkey, err := datastore.Put(repo.Scope, incompletekey, user)
	if err != nil {
		return 0, datastoreError(err)
	}
	user.ID = ID(userkey.IntID())

	mappingkey := datastore.NewKey(repo.Scope, "Credential", cred.ID, 0, nil)

	mapping := &credentialMapping{}
	mapping.UserKey = userkey
	mapping.Provider = cred.Provider
	mapping.Email = cred.Email
	mapping.Name = cred.Name

	_, err = datastore.Put(repo.Scope, mappingkey, mapping)
	if err != nil {
		return 0, datastoreError(err)
	}

	return user.ID, nil
}

func (repo *Datastore) ByCredentials(cred *auth.Credentials) (*User, error) {
	mapping := &credentialMapping{}

	mappingkey := datastore.NewKey(repo.Scope, "Credential", cred.ID, 0, nil)
	err := datastore.Get(repo.Scope, mappingkey, mapping)
	if err != nil {
		return nil, datastoreError(err)
	}

	user := &User{}
	user.ID = ID(mapping.UserKey.IntID())
	err = datastore.Get(repo.Scope, mapping.UserKey, user)

	return user, datastoreError(err)
}

func (repo *Datastore) ByID(id ID) (*User, error) {
	user := &User{}
	user.ID = id
	userkey := datastore.NewKey(repo.Scope, "User", "", int64(id), nil)
	err := datastore.Get(repo.Scope, userkey, user)
	return user, datastoreError(err)
}

func (repo *Datastore) Update(user *User) error {
	userkey := datastore.NewKey(repo.Scope, "User", "", int64(user.ID), nil)
	_, err := datastore.Put(repo.Scope, userkey, user)
	return datastoreError(err)
}
