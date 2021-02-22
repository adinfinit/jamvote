package postgresdb

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

func (db *DB) Migrate(ctx context.Context) error {
	_, err := db.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
		    id    INT8 not null,
		    name  TEXT not null,
		    email TEXT not null,
		    admin BOOL not null,

		    social JSONB not null,

		    PRIMARY KEY (id)
		);

 		CREATE TABLE IF NOT EXISTS credentials (
 		    provider TEXT not null,
 		    id       TEXT not null,
 		    email    TEXT not null,
 		    name     TEXT not null,

			user_id  INT8 not null,

 		    CONSTRAINT fk_credentials FOREIGN KEY (user_id) REFERENCES users(id)
 		    -- TODO add index 
 		);

		CREATE TABLE IF NOT EXISTS events (
		    id    TEXT not null,
		    name  TEXT not null,
		    theme TEXT not null,
		    info  TEXT not null,
		    
		    created_at TIMESTAMP not null
		);

		CREATE TABLE IF NOT EXISTS events_organizers (
		    event_id TEXT not null,
		    user_id  INT8 not null		    
		);

		CREATE TABLE IF NOT EXISTS events_jammers (
		    event_id TEXT not null,
		    user_id  INT8 not null
		);

		CREATE TABLE IF NOT EXISTS teams (
		    event_id text not null,
		    id       int8 not null,
		    
		    PRIMARY KEY (event_id, id)
		);
	`)
	return nil
}

func (db *DB) Events(ctx context.Context) event.Repo {
	return &Events{DB: db}
}

func (db *DB) Users(ctx context.Context) user.Repo {
	return &Users{DB: db}
}
