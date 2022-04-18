package pg

import (
	"database/sql"
	"fmt"

	"whisperd.io/whisperd/whisperd/db"
	"whisperd.io/whisperd/whisperd/db/store"
)

func NewProvider(opts db.Opts) (db.Provider, error) {
	p := &Provider{}
	return p, nil
}

type Provider struct {
}

func (p *Provider) DB() (*sql.DB, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p *Provider) Shouts(db *sql.DB) (store.Shouts, error) {
	return nil, fmt.Errorf("not implemented")
}
