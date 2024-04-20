package types

import (
	"github.com/shariqali-dev/discovery-trail/internal/models"
	"github.com/shariqali-dev/discovery-trail/internal/validator"
)

type CourseWithUnitTitles struct {
	UnitTitles []string
	models.Course
}

type CourseWithUnitsWithChapters struct {
	Units []UnitWithChapters
	models.Course
}

type UnitWithChapters struct {
	Chapters []models.Chapter
	models.Unit
}
type TemplateData struct {
	Nonce               string
	Flash               string
	IsAuthenticated     bool
	CSRFToken           string
	Account             models.Account
	Courses             []CourseWithUnitTitles
	Course              models.Course
	CourseUnitsChapters CourseWithUnitsWithChapters
	Form                map[string]any
}

type Unit struct {
	ID       uint8
	Name     string
	CourseID uint8
}

type CourseCreateForm struct {
	Title               string            `form:"title"`
	UnitCount           int               `form:"-"`
	UnitValues          map[string]string `form:"-"`
	validator.Validator `form:"-"`
}
