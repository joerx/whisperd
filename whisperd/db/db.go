package db

import (
	"database/sql"
	"errors"

	"whisperd.io/whisperd/whisperd/db/store"
)

var (
	ErrorNotFound      = errors.New("item not found")
	ErrorInvalidRecord = errors.New("invalid record")
)

type Opts struct {
	Driver         string
	SQLiteFileName string
}

// Provider interface for database objects. Implementations of this interface create database connection and storage
// implementations based on the selected database driver and possible other options.
type Provider interface {
	DB() (*sql.DB, error)
	Shouts(db *sql.DB) (store.Shouts, error)
}
