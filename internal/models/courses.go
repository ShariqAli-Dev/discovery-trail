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
	Processed bool
	AccountID string
}

type CourseModelInterface interface {
	Insert(name, image, account_id string) (string, error)
	All(account_id string) ([]Course, error)
	Get(course_id string) (Course, error)
	Process(courseID string) error
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

func (m *CourseModel) Get(course_id string) (Course, error) {
	sqlStatement := "SELECT id, name, image, processed, account_id FROM courses WHERE id = ?"
	row := m.DB.QueryRow(sqlStatement, course_id)

	var course Course
	err := row.Scan(&course.ID, &course.Name, &course.Image, &course.Processed, &course.AccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Course{}, ErrorNoRecord
		}
		return Course{}, err
	}
	return course, nil
}

func (m *CourseModel) All(account_id string) ([]Course, error) {
	sqlStatement := "SELECT id, name, image, processed FROM courses WHERE account_id = ?"
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
		err := rows.Scan(&course.ID, &course.Name, &course.Image, &course.Processed)
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

func (m *CourseModel) Process(courseID string) error {
	sqlStatement := "UPDATE courses SET processed = 1 WHERE id = ?"
	_, err := m.DB.Exec(sqlStatement, courseID)
	return err
}
