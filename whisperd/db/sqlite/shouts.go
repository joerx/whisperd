package sqlite

import (
	"context"
	"database/sql"
	"time"

	"whisperd.io/whisperd/whisperd"
	"whisperd.io/whisperd/whisperd/db"
	"whisperd.io/whisperd/whisperd/db/store"
)

const (
	isoFormat = "2006-01-02T15:04:05Z0700"
)

type sqliteStore struct {
	db *sql.DB
}

func NewShoutStore(db *sql.DB) store.Shouts {
	ss := &sqliteStore{db}
	return ss
}

func (ss *sqliteStore) GetAll(ctx context.Context) ([]whisperd.Shout, error) {
	rows, err := ss.db.QueryContext(ctx, "SELECT id, message, timestamp FROM shouts")
	if err != nil {
		return []whisperd.Shout{}, err
	}

	defer rows.Close()

	var (
		id        int64
		message   string
		timestamp string
	)

	sl := []whisperd.Shout{}

	for rows.Next() {
		if err := rows.Scan(&id, &message, &timestamp); err != nil {
			return sl, err
		}

		tt, err := time.Parse(isoFormat, timestamp)
		if err != nil {
			return sl, err
		}

		s := whisperd.Shout{ID: id, Message: message, Timestamp: tt}
		sl = append(sl, s)
	}

	return sl, nil
}

func (ss *sqliteStore) Get(ctx context.Context, id string) (whisperd.Shout, error) {
	return whisperd.Shout{}, nil
}

func (ss *sqliteStore) Insert(ctx context.Context, s whisperd.Shout) (whisperd.Shout, error) {
	if s.Message == "" {
		return s, db.ErrorInvalidRecord
	}

	tt := time.Now()

	res, err := ss.db.ExecContext(ctx,
		"INSERT INTO shouts (message, timestamp) VALUES ($1, $2)",
		s.Message,
		tt.Format(isoFormat),
	)
	if err != nil {
		return whisperd.Shout{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return whisperd.Shout{}, err
	}

	return whisperd.Shout{ID: id, Message: s.Message, Timestamp: tt}, nil
}

func (ss *sqliteStore) Delete(ctx context.Context, s whisperd.Shout) (whisperd.Shout, error) {
	return s, nil
}
