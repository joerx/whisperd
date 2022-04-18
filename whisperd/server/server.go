package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"whisperd.io/whisperd/whisperd/backend"
	"whisperd.io/whisperd/whisperd/db"
	"whisperd.io/whisperd/whisperd/frontend"
)

type Opts struct {
	Addr     string
	Roles    map[string]struct{}
	Frontend frontend.Opts
	DB       db.Opts
}

func New(opts Opts) (*http.Server, error) {
	be, err := backend.New(opts.DB)
	if err != nil {
		return nil, err
	}

	fe, err := frontend.New(frontend.Opts{})
	if err != nil {
		return nil, err
	}

	r := mux.NewRouter()
	r.PathPrefix("/api").Handler(be)
	r.PathPrefix("/").Handler(fe)

	s := &http.Server{
		Addr:    opts.Addr,
		Handler: r,
	}

	return s, nil
}
