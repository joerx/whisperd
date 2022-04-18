package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"whisperd.io/whisperd/whisperd/backend"
	"whisperd.io/whisperd/whisperd/db"
	"whisperd.io/whisperd/whisperd/frontend"
)

type Opts struct {
	Addr  string
	Roles map[string]struct{}
	DB    db.Opts
}

func New(opts Opts) (*http.Server, error) {
	r := mux.NewRouter()

	if _, ok := opts.Roles["backend"]; ok {
		be, err := backend.Handler(opts.DB)
		if err != nil {
			return nil, err
		}
		r.PathPrefix("/api").Handler(be)
	}

	if _, ok := opts.Roles["frontend"]; ok {
		fe, err := frontend.Handler()
		if err != nil {
			return nil, err
		}
		r.PathPrefix("/").Handler(fe)
	}

	s := &http.Server{
		Addr:    opts.Addr,
		Handler: r,
	}

	return s, nil
}
