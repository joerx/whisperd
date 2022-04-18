package sqlite

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"whisperd.io/whisperd/whisperd"
	"whisperd.io/whisperd/whisperd/db"
	"whisperd.io/whisperd/whisperd/db/store"
)

func initStore(t *testing.T) store.Shouts {
	db, err := Init(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	return NewShoutStore(db)
}

func withFixture(t *testing.T, s store.Shouts, rows []whisperd.Shout) []whisperd.Shout {
	out := make([]whisperd.Shout, 0, len(rows))
	for _, row := range rows {
		nr, err := s.Insert(context.Background(), row)
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, nr)
	}
	return out
}

func TestInsertShout(t *testing.T) {
	type test struct {
		input  whisperd.Shout
		err    error
		assert func(in, out whisperd.Shout) error
	}

	tests := []test{
		{
			input: whisperd.Shout{Message: "foo"},
			err:   nil,
			assert: func(si, so whisperd.Shout) error {
				if so.Message != si.Message {
					return fmt.Errorf("expected output message to be equal to input message '%s', but was '%s'", si.Message, so.Message)
				}
				if so.ID <= 0 {
					return fmt.Errorf("expected output record to have generated id greater than zero, but was %d", si.ID)
				}
				if so.Timestamp == si.Timestamp {
					return fmt.Errorf("expected output record to have updated timestamp but was %v", so.Timestamp)
				}
				if so.Timestamp.After(time.Now()) {
					return fmt.Errorf("expected output record to have timestamp not in the future, but was %v", so.Timestamp)
				}
				return nil
			},
		},
		{
			input:  whisperd.Shout{},
			err:    db.ErrorInvalidRecord,
			assert: func(si whisperd.Shout, so whisperd.Shout) error { return nil },
		},
	}

	ss := initStore(t)

	for _, tc := range tests {
		so, err := ss.Insert(context.Background(), tc.input)

		if err != tc.err {
			t.Fatalf("Expected error to be '%v' but got '%v'", tc.err, err)
		}
		if err := tc.assert(tc.input, so); err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetShouts(t *testing.T) {
	type test struct {
		shouts []whisperd.Shout
		err    error
		assert func(in, out []whisperd.Shout) error
	}

	tests := []test{
		{
			shouts: []whisperd.Shout{{Message: "foo"}, {Message: "bar"}, {Message: "baz"}},
			err:    nil,
			assert: func(s1, s2 []whisperd.Shout) error {
				for _, s := range s2 {
					if s.Message == "" {
						return errors.New("Expected record to have non-empty message")
					}
					if s.ID == 0 {
						return errors.New("Expected record to have non-zero id")
					}
					if (s.Timestamp == time.Time{}) {
						return errors.New("Expected record to have non-emtpy timestamp")
					}
				}
				return nil
			},
		},
	}

	for _, tc := range tests {
		ss := initStore(t)
		rows := withFixture(t, ss, tc.shouts)

		result, err := ss.GetAll(context.Background())

		if err != tc.err {
			t.Fatalf("Expected error to be '%v' but was '%v'", err, tc.err)
		}
		if len(result) != len(tc.shouts) {
			t.Fatalf("Expected result to have %d rows but got %d", len(tc.shouts), len(result))
		}
		if err := tc.assert(rows, result); err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetShout(t *testing.T) {
	shouts := []whisperd.Shout{{Message: "test get shout"}}
	ss := initStore(t)
	rows := withFixture(t, ss, shouts)

	type test struct {
		id     int64
		err    error
		assert func(id int64, s whisperd.Shout) error
	}

	tests := []test{
		{
			id:  1,
			err: nil,
			assert: func(id int64, s whisperd.Shout) error {
				if s.ID != id {
					return fmt.Errorf("Expected returned id to be '%d' but was '%d'", id, s.ID)
				}
				if s.Message != rows[0].Message {
					return fmt.Errorf("Expected returned message to be %s but was %s", rows[0].Message, s.Message)
				}
				if s.Timestamp.Format(isoFormat) != rows[0].Timestamp.Format(isoFormat) {
					return fmt.Errorf("Expected returned timestamp to be '%s' but was '%s'", rows[0].Timestamp, s.Timestamp)
				}
				return nil
			},
		},
		{
			id:     20,
			err:    db.ErrorNotFound,
			assert: func(id int64, s whisperd.Shout) error { return nil },
		},
	}

	for _, tc := range tests {
		result, err := ss.Get(context.Background(), tc.id)

		if err != tc.err {
			t.Fatalf("Expected error to be '%v' but got '%v'", tc.err, err)
		}

		if err := tc.assert(tc.id, result); err != nil {
			t.Fatal(err)
		}
	}
}
