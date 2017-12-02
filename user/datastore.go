package user

import (
	"context"
	"sort"

	"github.com/adinfinit/jamvote/auth"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

type Datastore struct {
	Context context.Context
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

func (repo *Datastore) Create(cred *auth.Credentials, user *User) (UserID, error) {
	incompletekey := datastore.NewIncompleteKey(repo.Context, "User", nil)
	userkey, err := datastore.Put(repo.Context, incompletekey, user)
	if err != nil {
		return 0, datastoreError(err)
	}
	user.ID = UserID(userkey.IntID())

	mappingkey := datastore.NewKey(repo.Context, "Credential", cred.ID, 0, nil)

	mapping := &credentialMapping{}
	mapping.UserKey = userkey
	mapping.Provider = cred.Provider
	mapping.Email = cred.Email
	mapping.Name = cred.Name

	_, err = datastore.Put(repo.Context, mappingkey, mapping)
	if err != nil {
		return 0, datastoreError(err)
	}

	return user.ID, nil
}

func (repo *Datastore) ByCredentials(cred *auth.Credentials) (*User, error) {
	if item, err := memcache.Get(repo.Context, "Credential_"+cred.ID); err == nil {
		user := &User{}
		if _, err := memcache.Gob.Get(repo.Context, "User_"+string(item.Value), user); err == nil {
			return user, nil
		}
	}

	mapping := &credentialMapping{}

	mappingkey := datastore.NewKey(repo.Context, "Credential", cred.ID, 0, nil)
	err := datastore.Get(repo.Context, mappingkey, mapping)
	if err != nil {
		return nil, datastoreError(err)
	}

	user := &User{}
	user.ID = UserID(mapping.UserKey.IntID())
	err = datastore.Get(repo.Context, mapping.UserKey, user)

	memcache.Add(repo.Context, &memcache.Item{
		Key:   "Credential_" + cred.ID,
		Value: []byte(UserID(mapping.UserKey.IntID()).String()),
	})

	memcache.Gob.Add(repo.Context, &memcache.Item{
		Key:    "User_" + UserID(mapping.UserKey.IntID()).String(),
		Object: user,
	})

	return user, datastoreError(err)
}

func (repo *Datastore) ByID(id UserID) (*User, error) {
	user := &User{}
	if _, err := memcache.Gob.Get(repo.Context, "User_"+id.String(), user); err == nil {
		return user, nil
	}

	user = &User{}
	user.ID = id
	userkey := datastore.NewKey(repo.Context, "User", "", int64(id), nil)
	err := datastore.Get(repo.Context, userkey, user)
	return user, datastoreError(err)
}

func (repo *Datastore) List() ([]*User, error) {
	users := []*User{}

	q := datastore.NewQuery("User")
	keys, err := q.GetAll(repo.Context, &users)
	for i, user := range users {
		user.ID = UserID(keys[i].IntID())
	}

	sort.Slice(users, func(i, k int) bool {
		return users[i].Name < users[k].Name
	})
	return users, datastoreError(err)
}

func (repo *Datastore) Update(user *User) error {
	userkey := datastore.NewKey(repo.Context, "User", "", int64(user.ID), nil)
	_, err := datastore.Put(repo.Context, userkey, user)
	if err == nil {
		memcache.Gob.Set(repo.Context, &memcache.Item{
			Key:    "User_" + user.ID.String(),
			Object: user,
		})
	}
	return datastoreError(err)
}
