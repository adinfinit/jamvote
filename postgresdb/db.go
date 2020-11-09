package datastoredb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/adinfinit/jamvote/event"
	"github.com/adinfinit/jamvote/user"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	db *sql.DB
}

func Open(ctx context.Context, dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("unable to open postgres db: %w", err)
	}
	return &DB{db}, nil
}

func (db *DB) Events(ctx context.Context) event.Repo {
	return &Events{DB: db}
}

func (db *DB) Users(ctx context.Context) user.Repo {
	return &Users{DB: db}
}
