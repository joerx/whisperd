package provider

import (
	"fmt"
	"strings"

	"whisperd.io/whisperd/whisperd/db"
	"whisperd.io/whisperd/whisperd/db/sqlite"
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
		return nil, fmt.Errorf("postgres driver is not implemented yet")
	default:
		supported := strings.Join([]string{SQLite, Postgres}, ", ")
		return nil, fmt.Errorf("invalid database driver %s, supported drivers are %s", opts.Driver, supported)
	}
}
