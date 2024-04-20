package models

import (
	"database/sql"
	"errors"
)

type Unit struct {
	ID       int64
	Name     string
	CourseID string
}
type UnitModelInterface interface {
	Insert(name string, course_id string) (int64, error)
	GetCourseUnitTitles(course_id string) ([]string, error)
	GetCourseUnits(course_id string) ([]Unit, error)
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
		return make([]string, 0), err
	}

	return titles, nil
}

func (m *UnitModel) GetCourseUnits(courseID string) ([]Unit, error) {
	sqlStatement := "SELECT id, name FROM units WHERE course_id = ?"
	rows, err := m.DB.Query(sqlStatement, courseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Unit{}, nil
		}
		return []Unit{}, err
	}

	var units []Unit
	for rows.Next() {
		var unit Unit
		unit.CourseID = courseID
		err := rows.Scan(&unit.ID, &unit.Name)
		if err != nil {
			return []Unit{}, err
		}
		units = append(units, unit)
	}
	if err = rows.Err(); err != nil {
		return []Unit{}, err
	}

	return units, nil
}
