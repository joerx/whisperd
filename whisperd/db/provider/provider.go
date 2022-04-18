package provider

import (
	"errors"

	"whisperd.io/whisperd/whisperd/db"
	"whisperd.io/whisperd/whisperd/db/pg"
	"whisperd.io/whisperd/whisperd/db/sqlite"
)

var (
	ErrorInvalidDriver = errors.New("invalid driver name")
)

const (
	SQLite   = "sqlite"
	Postgres = "postgres"
)

func New(opts db.Opts) (db.Provider, error) {
	switch opts.Driver {
	case SQLite:
		return sqlite.NewProvider(opts)
	case Postgres:
		return pg.NewProvider(opts)
	default:
		return nil, ErrorInvalidDriver
	}
}
