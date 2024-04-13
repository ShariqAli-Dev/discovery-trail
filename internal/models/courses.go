package models

import (
	"database/sql"
	"strings"

	"github.com/google/uuid"
)

type Course struct {
	ID    string
	Name  string
	Image string
}

type CourseModelInterface interface {
	Insert(name, image string) (string, error)
}

type CourseModel struct {
	DB *sql.DB
}

func (m *CourseModel) Insert(name, image string) (string, error) {
	id := uuid.New().String()
	sqlStatement := "INSERT INTO courses (id, name, image) VALUES(?, ?, ?)"
	_, err := m.DB.Exec(sqlStatement, id, name, strings.TrimSpace(image))

	return id, err
}
