package models

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/shariqali-dev/discovery-trail/internal/gpt"
)

type Question struct {
	ID        int64
	Question  string
	Answer    int64
	Options   []string
	ChapterID int64
}

type QuestionModelInterface interface {
	Insert(question gpt.GeneratedQuestion, chapterID int64) (int64, error)
	Delete(questionID int64) error
	GetFromChapterID(chapterID int64) (Question, error)
	ExistsFromChapterID(chapterID int64) (bool, error)
}

type QuestionModel struct {
	DB *sql.DB
}

func (m *QuestionModel) Insert(question gpt.GeneratedQuestion, chapterID int64) (int64, error) {
	optionsJSON, err := json.Marshal(question.Options)
	if err != nil {
		return 0, err
	}

	sqlStatement := "INSERT INTO questions (question, answer, options, chapter_id) VALUES (?, ?, ?, ?)"
	rowData, err := m.DB.Exec(sqlStatement, question.Question, question.Answer, optionsJSON, chapterID)
	if err != nil {
		return 0, err
	}
	return rowData.LastInsertId()
}

func (m *QuestionModel) Delete(questionID int64) error {
	sqlStatement := "DELETE FROM questions WHERE id = ?"
	_, err := m.DB.Exec(sqlStatement, questionID)
	return err
}

func (m *QuestionModel) GetFromChapterID(chapterID int64) (Question, error) {
	sqlStatement := "SELECT id, question, answer, options WHERE chapter_id = ?"
	var question Question
	question.ChapterID = chapterID
	if err := m.DB.QueryRow(sqlStatement, chapterID).Scan(&question.ID, &question.Question, &question.Answer, &question.Options); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return question, ErrorNoRecord
		}
		return question, err
	}

	return question, nil
}

func (m *QuestionModel) ExistsFromChapterID(chapterID int64) (bool, error) {
	var exists bool
	sqlStatement := "SELECT EXISTS(SELECT true FROM questions WHERE chapter_id = ?)"
	if err := m.DB.QueryRow(sqlStatement, chapterID).Scan(&exists); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return exists, err
		}
	}
	exists = true
	return exists, nil
}
