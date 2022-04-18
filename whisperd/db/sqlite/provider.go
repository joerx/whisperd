package sqlite

import (
	"database/sql"

	"whisperd.io/whisperd/whisperd/db"
	"whisperd.io/whisperd/whisperd/db/store"
)

func NewProvider(opts db.Opts) (db.Provider, error) {
	return &Provider{opts}, nil
}

type Provider struct {
	opts db.Opts
}

func (p *Provider) DB() (*sql.DB, error) {
	return Init(p.opts.SQLiteFileName)
}

func (p *Provider) Shouts(db *sql.DB) (store.Shouts, error) {
	return NewShoutStore(db), nil
}
