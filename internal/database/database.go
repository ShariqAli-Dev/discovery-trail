package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func OpenDB(dbURL, authToken string) (*sql.DB, error) {
	url := fmt.Sprintf("%s?authToken=%s", dbURL, authToken)
	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, fmt.Errorf("failed to open db %s: %s", url, err)
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func MigrateUp(db *sql.DB) error {
	setupScript, err := os.ReadFile("./internal/database/migrations/setup.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(setupScript))
	if err != nil {
		return err
	}
	return nil
}

func MigrateDown(db *sql.DB) error {
	setupScript, err := os.ReadFile("./internal/database/migrations/teardown.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(setupScript))
	if err != nil {
		return err
	}
	return nil
}
