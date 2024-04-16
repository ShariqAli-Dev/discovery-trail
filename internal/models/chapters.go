package models

import "database/sql"

type Chapter struct {
	ID                 int64
	Name               string
	YoutubeSearchQuery string
	VideoID            string
	Summary            string
	UnitID             int64
}

type ChapterModelInterface interface {
	Insert(name, youtubeQuery string, unitID int64) error
}

type ChapterModel struct {
	DB *sql.DB
}

func (m *ChapterModel) Insert(name, youtubeQuery string, unitID int64) error {
	sqlStatement := "INSERT INTO chapters (name, youtubeSearchQuery, unit_id) VALUES (?, ?, ?)"
	_, err := m.DB.Exec(sqlStatement, name, youtubeQuery, unitID)
	if err != nil {
		return err
	}
	return nil
}
