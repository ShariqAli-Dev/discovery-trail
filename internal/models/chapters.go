package models

import (
	"database/sql"
	"errors"
)

type Chapter struct {
	ID                 int64
	UnitID             int64
	Name               string
	YoutubeSearchQuery sql.NullString
	VideoID            sql.NullString
	Summary            sql.NullString
	QuestionsStatus    QuestionStatus
	QuestionAttempts   int64
}

var QuestionStatuses = map[QuestionStatusKey]QuestionStatus{
	Pending:   {Value: "Pending"},
	Error:     {Value: "Errored"},
	Completed: {Value: "Completed"},
}

const (
	Pending   QuestionStatusKey = "pending"
	Error     QuestionStatusKey = "error"
	Completed QuestionStatusKey = "completed"
)

type QuestionStatus struct {
	Value string
}

type QuestionStatusKey string

type ChapterModelInterface interface {
	Insert(name, youtubeQuery string, unitID int64) error
	GetUnitChapters(unitID int64) ([]Chapter, error)
	Get(chapterID int64) (Chapter, error)
	UpdateChapterQuestionStatus(chapterID int64, status QuestionStatus) error
}

type ChapterModel struct {
	DB *sql.DB
}

func (m *ChapterModel) Insert(name, youtubeQuery string, unitID int64) error {
	sqlStatement := "INSERT INTO chapters (name, youtubeSearchQuery, unit_id, questionsStatus,questionAttempts ) VALUES (?, ?, ?, ?, 0)"
	_, err := m.DB.Exec(sqlStatement, name, youtubeQuery, unitID, QuestionStatuses[Pending].Value)
	return err
}
func (m *ChapterModel) GetUnitChapters(unitID int64) ([]Chapter, error) {
	sqlStatement := "SELECT id, name, youtubeSearchQuery, videoID, summary, questionsStatus, questionAttempts FROM chapters WHERE unit_id = ?"
	rows, err := m.DB.Query(sqlStatement, unitID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Chapter{}, nil
		}
		return []Chapter{}, err
	}

	var chapters []Chapter
	for rows.Next() {
		var chapter Chapter
		chapter.UnitID = int64(unitID)
		err := rows.Scan(&chapter.ID, &chapter.Name, &chapter.YoutubeSearchQuery, &chapter.VideoID, &chapter.Summary, &chapter.QuestionsStatus.Value, &chapter.QuestionAttempts)
		if err != nil {
			return []Chapter{}, err
		}
		chapters = append(chapters, chapter)
	}
	if err = rows.Err(); err != nil {
		return []Chapter{}, err
	}
	return chapters, nil
}

func (m *ChapterModel) UpdateChapterQuestionStatus(chapterID int64, status QuestionStatus) error {
	sqlStatement := "UPDATE chapters SET questionsStatus = ? WHERE id = ?"
	_, err := m.DB.Exec(sqlStatement, status.Value, chapterID)
	return err
}
func (m *ChapterModel) Get(chapterID int64) (Chapter, error) {
	sqlStatement := "SELECT unit_id, name, youtubeSearchQuery, videoID, summary, questionsStatus, questionAttempts FROM chapters where id = ?"
	var chapter Chapter
	chapter.ID = chapterID
	if err := m.DB.QueryRow(sqlStatement, chapterID).Scan(&chapter.UnitID, &chapter.Name, &chapter.YoutubeSearchQuery, &chapter.VideoID, &chapter.Summary, &chapter.QuestionsStatus.Value, &chapter.QuestionAttempts); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return chapter, ErrorNoRecord
		}
		return chapter, err
	}
	return chapter, nil
}
