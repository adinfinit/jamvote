package datastoredb

import (
	"context"
	"sort"

	"cloud.google.com/go/datastore"

	"github.com/adinfinit/jamvote/auth"
	"github.com/adinfinit/jamvote/internal/natural"
	"github.com/adinfinit/jamvote/user"
)

// Users implements user.Repo.
type Users struct {
	Context context.Context
	Client  *datastore.Client
}

// credentialMapping is info that is stored in the datastore.
type credentialMapping struct {
	UserKey  *datastore.Key
	Provider string `datastore:",noindex"`
	Email    string `datastore:",noindex"`
	Name     string `datastore:",noindex"`
}

// usersError converts datastore error to a domain error.
func usersError(err error) error {
	if err == datastore.ErrNoSuchEntity {
		return user.ErrNotExists
	}
	return err
}

// Create creates a new user with the specified credentials and user info.
func (repo *Users) Create(cred *auth.Credentials, u *user.User) (user.UserID, error) {
	incompletekey := datastore.IncompleteKey("User", nil)
	userkey, err := repo.Client.Put(repo.Context, incompletekey, u)
	if err != nil {
		return 0, usersError(err)
	}
	u.ID = user.UserID(userkey.ID)

	mappingkey := datastore.NameKey("Credential", cred.ID, nil)

	mapping := &credentialMapping{}
	mapping.UserKey = userkey
	mapping.Provider = cred.Provider
	mapping.Email = cred.Email
	mapping.Name = cred.Name

	_, err = repo.Client.Put(repo.Context, mappingkey, mapping)
	if err != nil {
		return 0, usersError(err)
	}

	return u.ID, nil
}

// ByCredentials finds user info based on credentials.
func (repo *Users) ByCredentials(cred *auth.Credentials) (*user.User, error) {
	u := &user.User{}
	if appCache.Get("Credential_"+cred.ID, u) {
		return u, nil
	}

	mapping := &credentialMapping{}

	mappingkey := datastore.NameKey("Credential", cred.ID, nil)
	err := repo.Client.Get(repo.Context, mappingkey, mapping)
	if err != nil {
		return nil, usersError(err)
	}

	u = &user.User{}
	u.ID = user.UserID(mapping.UserKey.ID)
	err = repo.Client.Get(repo.Context, mapping.UserKey, u)

	if err == nil {
		appCache.Set("Credential_"+cred.ID, u)
		appCache.Set("User_"+u.ID.String(), u)
	}

	return u, usersError(err)
}

// ByID returns user by ID.
func (repo *Users) ByID(id user.UserID) (*user.User, error) {
	u := &user.User{}
	if appCache.Get("User_"+id.String(), u) {
		return u, nil
	}

	u = &user.User{}
	u.ID = id
	userkey := datastore.IDKey("User", int64(id), nil)
	err := repo.Client.Get(repo.Context, userkey, u)
	return u, usersError(err)
}

// List returns all users.
func (repo *Users) List() ([]*user.User, error) {
	users := []*user.User{}

	q := datastore.NewQuery("User")
	keys, err := repo.Client.GetAll(repo.Context, q, &users)
	for i, u := range users {
		u.ID = user.UserID(keys[i].ID)
	}

	sort.Slice(users, func(i, k int) bool {
		return natural.Less(users[i].Name, users[k].Name)
	})
	return users, usersError(err)
}

// FindCredentialByEmail scans all credentials for a matching email and returns the associated UserID.
func (repo *Users) FindCredentialByEmail(email string) (user.UserID, error) {
	var mappings []credentialMapping
	_, err := repo.Client.GetAll(repo.Context, datastore.NewQuery("Credential"), &mappings)
	if err != nil {
		return 0, err
	}

	for _, m := range mappings {
		if m.Email == email {
			return user.UserID(m.UserKey.ID), nil
		}
	}

	return 0, user.ErrNotExists
}

// CreateCredentialAlias creates a new credential mapping pointing to an existing user.
func (repo *Users) CreateCredentialAlias(cred *auth.Credentials, existingUserID user.UserID) error {
	userkey := datastore.IDKey("User", int64(existingUserID), nil)
	mappingkey := datastore.NameKey("Credential", cred.ID, nil)

	mapping := &credentialMapping{
		UserKey:  userkey,
		Provider: cred.Provider,
		Email:    cred.Email,
		Name:     cred.Name,
	}

	_, err := repo.Client.Put(repo.Context, mappingkey, mapping)
	return err
}

// Update updates a user.
func (repo *Users) Update(u *user.User) error {
	userkey := datastore.IDKey("User", int64(u.ID), nil)
	_, err := repo.Client.Put(repo.Context, userkey, u)
	if err == nil {
		appCache.Set("User_"+u.ID.String(), u)
	}
	return usersError(err)
}
