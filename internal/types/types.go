package types

import (
	"github.com/shariqali-dev/discovery-trail/internal/models"
	"github.com/shariqali-dev/discovery-trail/internal/validator"
)

type TemplateCourseWithUnitTitles struct {
	UnitTitles []string
	models.Course
}

type TemplateData struct {
	Nonce           string
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
	Account         models.Account
	Courses         []TemplateCourseWithUnitTitles
	Form            map[string]any
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
