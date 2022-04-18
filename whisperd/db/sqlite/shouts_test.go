package sqlite

import (
	"context"
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

func TestInsertShout(t *testing.T) {
	ss := initStore(t)

	si := whisperd.Shout{Message: "foo"}
	so, err := ss.Insert(context.Background(), si)
	if err != nil {
		t.Fatal(err)
	}

	if so.Message != si.Message {
		t.Fatalf("expected output message to be equal to input message '%s', but was '%s'", si.Message, so.Message)
	}
	if so.ID == 0 {
		t.Fatalf("expected output record to have generated id, but id was %d", si.ID)
	}
	if so.ID < 0 {
		t.Fatalf("expected output record to have generated id greated than zero, but was %d", si.ID)
	}
	if so.Timestamp == si.Timestamp {
		t.Fatalf("expected output record to have updated timestamp but was %v", so.Timestamp)
	}
	if so.Timestamp.After(time.Now()) {
		t.Fatalf("expected output record to have timestamp not in the future, but was %v", so.Timestamp)
	}
}

func TestInsertShoutErrors(t *testing.T) {
	ss := initStore(t)

	si := whisperd.Shout{}
	_, err := ss.Insert(context.Background(), si)

	if err != db.ErrorInvalidRecord {
		t.Fatalf("expected error '%v' but got '%v'", db.ErrorInvalidRecord, err)
	}
}
