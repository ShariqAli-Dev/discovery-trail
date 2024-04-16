package models

import (
	"database/sql"
	"errors"
)

type Unit struct {
	ID       string
	Name     string
	CourseID int
}
type UnitModelInterface interface {
	Insert(name string, course_id string) (int64, error)
	GetCourseUnitTitles(course_id string) ([]string, error)
}

type UnitModel struct {
	DB *sql.DB
}

func (m *UnitModel) Insert(name string, course_id string) (int64, error) {
	sqlStatement := "INSERT INTO units (name, course_id) VALUES(?, ?)"
	res, err := m.DB.Exec(sqlStatement, name, course_id)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (m *UnitModel) GetCourseUnitTitles(course_id string) ([]string, error) {
	sqlStatement := "SELECT name FROM units WHERE course_id = ?"
	rows, err := m.DB.Query(sqlStatement, course_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]string, 0), nil
		}
		return make([]string, 0), err
	}

	var titles []string
	for rows.Next() {
		var title string
		err := rows.Scan(&title)
		if err != nil {
			return make([]string, 0), nil
		}
		titles = append(titles, title)
	}

	if err = rows.Err(); err != nil {
		return make([]string, 0), nil
	}

	return titles, nil
}
