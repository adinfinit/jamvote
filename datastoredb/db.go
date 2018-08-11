package datastoredb

import (
	"context"

	"github.com/adinfinit/jamvote/event"
	"github.com/adinfinit/jamvote/user"
)

type DB struct{}

func (db *DB) Events(ctx context.Context) event.Repo {
	return &Events{ctx}
}

func (db *DB) Users(ctx context.Context) user.Repo {
	return &Users{ctx}
}
