package datastoredb

import (
	"context"

	"github.com/adinfinit/jamvote/event"
	"github.com/adinfinit/jamvote/user"
)

// DB implements master database.
type DB struct{}

// Events returns event.Repo.
func (db *DB) Events(ctx context.Context) event.Repo {
	return &Events{ctx}
}

// Users returns user.Repo.
func (db *DB) Users(ctx context.Context) user.Repo {
	return &Users{ctx}
}
