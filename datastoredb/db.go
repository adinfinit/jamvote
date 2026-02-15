package datastoredb

import (
	"context"

	"cloud.google.com/go/datastore"

	"github.com/adinfinit/jamvote/event"
	"github.com/adinfinit/jamvote/user"
)

// DB implements master database.
type DB struct {
	Client *datastore.Client
}

// Events returns event.Repo.
func (db *DB) Events(ctx context.Context) event.Repo {
	return &Events{Context: ctx, Client: db.Client}
}

// Users returns user.Repo.
func (db *DB) Users(ctx context.Context) user.Repo {
	return &Users{Context: ctx, Client: db.Client}
}
