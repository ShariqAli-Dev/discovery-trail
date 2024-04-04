package models

import (
	"database/sql"
	"errors"
	"strings"
)

type User struct {
	ID    string
	Name  string
	Email string
}

type UserModelInterface interface {
	Insert(id, name, email string) error
	Exists(id string) (bool, error)
	Get(id string) (User, error)
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(id, name, email string) error {
	sqlStatement := "INSERT INTO users (id, name, email) VALUES(?, ?, ?)"
	_, err := m.DB.Exec(sqlStatement, id, name, email)

	return err

}

func (m *UserModel) Exists(id string) (bool, error) {
	var exists bool
	sqlStatement := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := m.DB.QueryRow(sqlStatement, strings.Trim(id, `"`)).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrorNoRecord
		} else {
			return false, err
		}
	}

	return exists, nil
}

func (m *UserModel) Get(id string) (User, error) {
	var user User
	sqlStatement := "SLELECT id, name, email FROM users WHERE id = ?"

	err := m.DB.QueryRow(sqlStatement, strings.Trim(id, `"`)).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrorNoRecord
		} else {
			return User{}, err
		}
	}

	return user, nil
}
