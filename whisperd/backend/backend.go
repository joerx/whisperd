package backend

import (
	"net/http"

	"whisperd.io/whisperd/whisperd/db"
	"whisperd.io/whisperd/whisperd/db/provider"
)

func Handler(dbOpts db.Opts) (http.Handler, error) {
	p, err := provider.New(dbOpts)
	if err != nil {
		return nil, err
	}

	conn, err := p.DB()
	if err != nil {
		return nil, err
	}

	ss, err := p.Shouts(conn)
	if err != nil {
		return nil, err
	}

	b := &handler{shouts: ss}
	return newRouter(b), nil
}
