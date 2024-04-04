package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
)

type SessionModelInterface interface {
	Create(data []byte) (string, error)
	Exists(token string) (bool, error)
	Destroy(token string) error
	GetUserID(token string) (string, error)
}

type Session struct {
	Token  string
	Data   []byte
	Expiry time.Time
}

type SessionModel struct {
	DB *sql.DB
}

func (m *SessionModel) Create(data []byte) (string, error) {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	sessionID := hex.EncodeToString(randomBytes)
	daysUntilSessionExpires := 7
	expiryTime := time.Now().UTC().Add(time.Duration(daysUntilSessionExpires) * 24 * time.Hour)
	expiryTimeString := expiryTime.Format("2006-01-02 15:04:05")

	sqlStatement := "INSERT INTO sessions (token, data, expiry) VALUES (?, ?, ?)"
	_, err = m.DB.Exec(sqlStatement, sessionID, data, expiryTimeString)

	return sessionID, err
}

func (m *SessionModel) Exists(token string) (bool, error) {
	var exists bool
	sqlStatement := "SELECT EXISTS(SELECT true FROM SESSIONS where token = ?)"
	err := m.DB.QueryRow(sqlStatement, token).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrorNoRecord
		} else {
			return false, err
		}
	}

	return exists, nil
}

func (m *SessionModel) Destroy(token string) error {
	sqlStatement := "DELETE FROM sessions WHERE token = ?"
	_, err := m.DB.Exec(sqlStatement, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorNoRecord
		} else {
			return err
		}
	}

	return nil
}

func (m *SessionModel) GetUserID(token string) (string, error) {
	var userID string
	sqlStatement := "SELECT data FROM sessions WHERE token = ?"
	err := m.DB.QueryRow(sqlStatement, token).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userID, ErrorNoRecord
		} else {
			return userID, err
		}
	}
	return userID, nil
}
