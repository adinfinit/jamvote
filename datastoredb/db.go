package datastoredb

import (
	"context"

	"github.com/adinfinit/jamvote/event"
)

type DB struct{}

func (db *DB) Events(ctx context.Context) event.Repo {
	return &Events{ctx}
}
