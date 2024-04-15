package models

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type Course struct {
	ID        string
	Name      string
	Image     string
	AccountID string
}

type CourseModelInterface interface {
	Insert(name, image, account_id string) (string, error)
	All(account_id string) ([]Course, error)
}

type CourseModel struct {
	DB *sql.DB
}

func (m *CourseModel) Insert(name, image, account_id string) (string, error) {
	id := uuid.New().String()

	sqlStatement := "INSERT INTO courses (id, name, image, account_id) VALUES(?, ?, ?, ?)"
	_, err := m.DB.Exec(sqlStatement, id, name, image, account_id)

	return id, err
}

func (m *CourseModel) All(account_id string) ([]Course, error) {
	sqlStatement := "SELECT id, name, image FROM courses WHERE account_id = ?"
	rows, err := m.DB.Query(sqlStatement, account_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Course{}, nil
		}
		return []Course{}, err
	}

	var courses []Course
	for rows.Next() {
		var course Course
		err := rows.Scan(&course.ID, &course.Name, &course.Image)
		if err != nil {
			return []Course{}, err
		}
		courses = append(courses, course)
	}
	if err = rows.Err(); err != nil {
		return []Course{}, nil
	}
	return courses, nil
}
