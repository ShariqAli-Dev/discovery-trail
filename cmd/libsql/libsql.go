package libsql

import (
	"database/sql"
	"time"

	"github.com/alexedwards/scs/v2"
)

type libsqlStore struct {
	db *sql.DB
}

func (s *libsqlStore) Delete(token string) error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE token = ?", token)
	return err
}

func (s *libsqlStore) Find(token string) ([]byte, bool, error) {
	var data []byte
	var expiry time.Time

	err := s.db.QueryRow("SELECT data, expiry FROM sessions WHERE token = ?", token).Scan(&data, &expiry)
	switch {
	case err == sql.ErrNoRows:
		return nil, false, nil // session not found
	case err != nil:
		return nil, false, err // error occurred
	default:
		return data, true, nil // session found
	}
}

func (s *libsqlStore) Commit(token string, data []byte, expiry time.Time) error {
	_, err := s.db.Exec("INSERT INTO sessions (token, data, expiry) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE data = VALUES(data), expiry = VALUES(expiry)", token, data, expiry)
	return err
}

func NewLibSQLStore(db *sql.DB) scs.Store {
	return &libsqlStore{db: db}
}
