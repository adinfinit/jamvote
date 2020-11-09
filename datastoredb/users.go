package datastoredb

import (
	"context"
	"errors"
	"sort"

	"github.com/adinfinit/jamvote/auth"
	"github.com/adinfinit/jamvote/internal/natural"
	"github.com/adinfinit/jamvote/user"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

type Users struct {
	Context context.Context
}

type credentialMapping struct {
	UserKey  *datastore.Key
	Provider string `datastore:",noindex"`
	Email    string `datastore:",noindex"`
	Name     string `datastore:",noindex"`
}

func usersError(err error) error {
	if errors.Is(err, datastore.ErrNoSuchEntity) {
		return user.ErrNotExists
	}
	return err
}

func (repo *Users) Create(cred *auth.Credentials, u *user.User) (user.UserID, error) {
	incompletekey := datastore.NewIncompleteKey(repo.Context, "User", nil)
	userkey, err := datastore.Put(repo.Context, incompletekey, u)
	if err != nil {
		return 0, usersError(err)
	}
	u.ID = user.UserID(userkey.IntID())

	mappingkey := datastore.NewKey(repo.Context, "Credential", cred.ID, 0, nil)

	mapping := &credentialMapping{}
	mapping.UserKey = userkey
	mapping.Provider = cred.Provider
	mapping.Email = cred.Email
	mapping.Name = cred.Name

	_, err = datastore.Put(repo.Context, mappingkey, mapping)
	if err != nil {
		return 0, usersError(err)
	}

	return u.ID, nil
}

func (repo *Users) ByCredentials(cred *auth.Credentials) (*user.User, error) {
	if item, err := memcache.Get(repo.Context, "Credential_"+cred.ID); err == nil {
		u := &user.User{}
		if _, err := memcache.Gob.Get(repo.Context, "User_"+string(item.Value), u); err == nil {
			return u, nil
		}
	}

	mapping := &credentialMapping{}

	mappingkey := datastore.NewKey(repo.Context, "Credential", cred.ID, 0, nil)
	err := datastore.Get(repo.Context, mappingkey, mapping)
	if err != nil {
		return nil, usersError(err)
	}

	u := &user.User{}
	u.ID = user.UserID(mapping.UserKey.IntID())
	err = datastore.Get(repo.Context, mapping.UserKey, u)

	memcache.Add(repo.Context, &memcache.Item{
		Key:   "Credential_" + cred.ID,
		Value: []byte(user.UserID(mapping.UserKey.IntID()).String()),
	})

	memcache.Gob.Add(repo.Context, &memcache.Item{
		Key:    "User_" + user.UserID(mapping.UserKey.IntID()).String(),
		Object: u,
	})

	return u, usersError(err)
}

func (repo *Users) ByID(id user.UserID) (*user.User, error) {
	u := &user.User{}
	if _, err := memcache.Gob.Get(repo.Context, "User_"+id.String(), u); err == nil {
		return u, nil
	}

	u = &user.User{}
	u.ID = id
	userkey := datastore.NewKey(repo.Context, "User", "", int64(id), nil)
	err := datastore.Get(repo.Context, userkey, u)
	return u, usersError(err)
}

func (repo *Users) List() ([]*user.User, error) {
	users := []*user.User{}

	q := datastore.NewQuery("User")
	keys, err := q.GetAll(repo.Context, &users)
	for i, u := range users {
		u.ID = user.UserID(keys[i].IntID())
	}

	sort.Slice(users, func(i, k int) bool {
		return natural.Less(users[i].Name, users[k].Name)
	})
	return users, usersError(err)
}

func (repo *Users) Update(u *user.User) error {
	userkey := datastore.NewKey(repo.Context, "User", "", int64(u.ID), nil)
	_, err := datastore.Put(repo.Context, userkey, u)
	if err == nil {
		memcache.Gob.Set(repo.Context, &memcache.Item{
			Key:    "User_" + u.ID.String(),
			Object: u,
		})
	}
	return usersError(err)
}
