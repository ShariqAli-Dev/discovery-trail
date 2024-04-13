package models

import (
	"database/sql"
	"errors"
)

type Account struct {
	ID      string
	Name    string
	Email   string
	Credits int8
}

type AccountModelInterface interface {
	Insert(id, name, email string) error
	Exists(id string) (bool, error)
	Get(id string) (Account, error)
}

type AccountModel struct {
	DB *sql.DB
}

func (m *AccountModel) Insert(id, name, email string) error {
	sqlStatement := "INSERT INTO accounts (id, name, email, credits) VALUES(?, ?, ?, 10)"
	_, err := m.DB.Exec(sqlStatement, id, name, email)

	return err

}

func (m *AccountModel) Exists(id string) (bool, error) {
	var exists bool
	sqlStatement := "SELECT EXISTS(SELECT true FROM accounts WHERE id = ?)"
	err := m.DB.QueryRow(sqlStatement, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrorNoRecord
		} else {
			return false, err
		}
	}

	return exists, nil
}

func (m *AccountModel) Get(id string) (Account, error) {
	var account Account
	sqlStatement := "SELECT id, name, email, credits FROM accounts WHERE id = ?"
	err := m.DB.QueryRow(sqlStatement, id).Scan(&account.ID, &account.Name, &account.Email, &account.Credits)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Account{}, ErrorNoRecord
		} else {
			return Account{}, err
		}
	}

	return account, nil
}
