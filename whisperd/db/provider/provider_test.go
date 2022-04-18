package provider

import (
	"reflect"
	"testing"

	"whisperd.io/whisperd/whisperd/db"
	"whisperd.io/whisperd/whisperd/db/pg"
	"whisperd.io/whisperd/whisperd/db/sqlite"
)

func TestProvider(t *testing.T) {
	type test struct {
		opts  db.Opts
		err   error
		proto db.Provider
	}

	tests := []test{
		{
			opts:  db.Opts{Driver: SQLite},
			err:   nil,
			proto: &sqlite.Provider{},
		},
		{
			opts:  db.Opts{Driver: Postgres},
			err:   nil,
			proto: &pg.Provider{},
		},
		{
			opts:  db.Opts{},
			err:   ErrorInvalidDriver,
			proto: nil,
		},
	}

	for _, tc := range tests {
		p, err := New(tc.opts)

		if err != tc.err {
			t.Fatalf("Expected error to be '%v' but was '%v'", tc.err, err)
		}

		eType := reflect.TypeOf(tc.proto)
		aType := reflect.TypeOf(p)
		if eType != aType {
			t.Fatalf("expected provider of type '%s', but was '%s'", eType, aType)
		}
	}
}
