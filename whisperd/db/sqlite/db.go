package sqlite

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Init(filename string) (*sql.DB, error) {
	log.Printf("Connecting to sqlite database, filename %s", filename)

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	if err := initSchema(db); err != nil {
		return nil, err
	}

	return db, nil
}

func initSchema(db *sql.DB) error {
	log.Println("Initializing database schema")

	query := `
	CREATE TABLE IF NOT EXISTS shouts(
		id INTEGER PRIMARY KEY ASC AUTOINCREMENT,
		message TEXT NOT NULL,
		timestamp TEXT NOT NULL
	)
	`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}
