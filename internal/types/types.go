package types

import (
	"github.com/shariqali-dev/discovery-trail/internal/models"
	"github.com/shariqali-dev/discovery-trail/internal/validator"
)

type TemplateData struct {
	Nonce           string
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
	Account         models.Account
	Courses         []models.Course
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
