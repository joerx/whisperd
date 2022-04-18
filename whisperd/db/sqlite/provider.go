package sqlite

import (
	"database/sql"

	"whisperd.io/whisperd/whisperd/db"
	"whisperd.io/whisperd/whisperd/db/store"
)

func NewProvider(opts db.Opts) (db.Provider, error) {
	return &provider{opts}, nil
}

type provider struct {
	opts db.Opts
}

func (p *provider) DB() (*sql.DB, error) {
	return Init(p.opts.SQLiteFileName)
}

func (p *provider) Shouts(db *sql.DB) (store.Shouts, error) {
	return NewShoutStore(db), nil
}
